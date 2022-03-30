package logstf

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const timeout = 5 * time.Minute

// Due to the fact log's playedAt filed points to gameserver time,
// and it can be non-UTC timezone, hours difference must be applied
const hourShift = 36

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

func (c Client) SearchLogs(players, maps []string, playedAt time.Time) ([]Log, error) {
	resp, err := c.getLogsWithPlayers(players)
	if err != nil {
		return nil, fmt.Errorf("failed to get players from logs.tf api: %v", err)
	}
	logs := filterLogs(maps, playedAt, resp.Logs)

	return logs, nil
}

func filterLogs(maps []string, playedAt time.Time, logs []Log) []Log {
	mapsWhitelist := make(map[string]struct{})

	validLogs := make([]Log, 0)

	for _, m := range maps {
		mapsWhitelist[m] = struct{}{}
	}

	dayBeforePlayed := playedAt.Add(-hourShift * time.Hour)
	dayAfterPlayed := playedAt.Add(hourShift * time.Hour)

	for _, log := range logs {
		if _, ok := mapsWhitelist[log.Map]; ok {
			logPlayedAt := time.Unix(int64(log.Date), 0)

			if logPlayedAt.After(dayBeforePlayed) && logPlayedAt.Before(dayAfterPlayed) {
				validLogs = append(validLogs, log)
			}
		}
	}
	return validLogs
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
