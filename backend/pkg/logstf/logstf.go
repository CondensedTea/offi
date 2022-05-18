package logstf

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	timeout           = 5 * time.Minute
	matchPlayedOffset = 3 * 24 * time.Hour
)

type Getter interface {
	Get(string) (*http.Response, error)
}

type Client struct {
	httpClient Getter
}

func New() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: timeout},
	}
}

func (c Client) SearchLogs(players, maps []string, playedAt time.Time) ([]Log, []Log, error) {
	started := time.Now()
	defer func() {
		logsTfSearchTime.WithLabelValues().Observe(time.Since(started).Seconds())
	}()

	resp, err := c.getLogsWithPlayers(players)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get players from logs.tf api: %v", err)
	}
	matchLogs, secondaryLogs := filterLogs(maps, resp.Logs, playedAt)

	return matchLogs, secondaryLogs, nil
}

func (c Client) GetLog(id int) (Log, error) {
	u := fmt.Sprintf("https://logs.tf/api/v1/log/%d", id)
	resp, err := c.httpClient.Get(u)
	if err != nil {
		return Log{}, err
	}
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return Log{}, fmt.Errorf("api returned bad staus: %d; %s", resp.StatusCode, string(b))
	}
	defer resp.Body.Close()

	var r FullLog
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return Log{}, err
	}
	return r.Info, nil
}

// getLogsWithPlayers gets logs with given players from logs.tf API
func (c Client) getLogsWithPlayers(players []string) (*Response, error) {
	query := "player=" + strings.Join(players, ",")

	u := fmt.Sprintf("https://logs.tf/api/v1/log?%s", query)
	resp, err := c.httpClient.Get(u)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("api returned bad staus: %d; %s", resp.StatusCode, string(b))
	}
	defer resp.Body.Close()
	var r Response

	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func filterLogs(maps []string, logs []Log, playedAt time.Time) (matchLogs, combinedLogs []Log) {
	for _, log := range logs {
		primary, valid := matchIsPrimary(playedAt, log.Date, log.Map, maps)
		if !valid {
			continue
		}
		if !primary {
			combinedLogs = append(combinedLogs, log)
		} else {
			matchLogs = append(matchLogs, log)
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
			logrus.Errorf("etf2l returned map without pattern [gamemode]_[mapname]: %s", m)
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
	return strings.Join(logMapParts[:genericMapItemLength], "_")
}
