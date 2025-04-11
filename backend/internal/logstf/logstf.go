package logstf

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	info "offi/internal/build_info"
	"offi/internal/tracing"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/transport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var logger = slog.With("component", "logs-tf")

const (
	matchPlayedOffset = 3 * 24 * time.Hour
)

type Client struct {
	client *http.Client
	tracer trace.Tracer
}

func NewClient() *Client {
	return &Client{
		client: &http.Client{
			Transport: transport.Chain(
				http.DefaultTransport,
				tracing.OTelHTTPTransport,
				transport.SetHeader("User-Agent", "offi-backend/"+info.Version),
			),
		},
		tracer: otel.GetTracerProvider().Tracer("logstf"),
	}
}

func (c *Client) SearchLogs(ctx context.Context, players []int, maps []string, playedAt time.Time) ([]Log, []Log, error) {
	ctx, span := c.tracer.Start(ctx, "logstf.SearchLogs")
	defer span.End()

	span.SetAttributes(attribute.IntSlice("player_ids", players))

	resp, err := c.getLogsWithPlayers(ctx, players)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get players from logs.tf api: %v", err)
	}
	matchLogs, secondaryLogs := filterLogs(maps, resp.Logs, playedAt)

	return matchLogs, secondaryLogs, nil
}

// getLogsWithPlayers gets logs with given players from logs.tf API
func (c *Client) getLogsWithPlayers(ctx context.Context, players []int) (*Response, error) {
	var b strings.Builder
	for _, steamID := range players {
		if b.Len() > 0 {
			b.WriteString(",")
		}
		b.WriteString(strconv.Itoa(steamID))
	}

	// TODO: use undocumented format parameter:
	// https://github.com/alevoska/logstf-web/blob/master/pylogstf/controllers/api.py#L156-L165

	u := fmt.Sprintf("https://logs.tf/api/v1/log?player=%s", b.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("creating request: %v", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("doing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("logs.tf API returned non-200 status: %d", resp.StatusCode)
	}

	var r Response
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("decoding response: %v", err)
	}
	return &r, nil
}

func filterLogs(maps []string, logs []Log, playedAt time.Time) (matchLogs, combinedLogs []Log) {
	// TODO: try to use ETF2L api-v2 data about map order, GC status, map scores

	for _, log := range logs {
		primary, valid := matchIsPrimary(playedAt, log.Date, log.Map, maps)
		if !valid {
			continue
		}
		if primary {
			matchLogs = append(matchLogs, log)
		} else {
			combinedLogs = append(combinedLogs, log)
		}
	}
	return matchLogs, combinedLogs
}

func matchIsPrimary(matchDate time.Time, logDate int64, logMap string, maps []string) (primary, valid bool) {
	matchDateMinusOffset := matchDate.Add(-matchPlayedOffset)
	matchDatePlusOffset := matchDate.Add(matchPlayedOffset)
	logPlayedAt := time.Unix(logDate, 0)

	if logPlayedAt.Before(matchDateMinusOffset) || logPlayedAt.After(matchDatePlusOffset) {
		return false, false
	}

	if mapIsNotValid(maps, logMap) {
		return false, true
	}
	return true, true
}

func mapIsNotValid(maps []string, logMap string) bool {
	mapsWhitelist := make(map[string]struct{})

	for _, m := range maps {
		genericMap := getGenericMapName(m)
		if genericMap == "" {
			logger.Error("etf2l returned map without pattern [gamemode]_[mapname]", "map_name", m)
			return true
		}
		mapsWhitelist[genericMap] = struct{}{}
	}

	genericLogMap := getGenericMapName(logMap)
	if genericLogMap == "" {
		return true
	}
	if _, ok := mapsWhitelist[genericLogMap]; ok {
		return false
	}
	return true
}

func getGenericMapName(mapName string) string {
	const genericMapItemLength = 2

	logMapParts := strings.Split(mapName, "_")
	if len(logMapParts) < 2 {
		return ""
	}
	return strings.ToLower(strings.Join(logMapParts[:genericMapItemLength], "_"))
}
