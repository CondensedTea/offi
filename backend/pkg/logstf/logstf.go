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

func filterLogs(maps []string, logs []Log, playedAt time.Time) (matchLogs, combinedLogs []Log) {
	mapsWhitelist := make(map[string]struct{})

	for _, m := range maps {
		mapsWhitelist[m] = struct{}{}
	}

	for _, log := range logs {

		matchPlayedAtMinusOffset := playedAt.Add(-matchPlayedOffset)
		matchPlayedAtPlusOffset := playedAt.Add(matchPlayedOffset)

		logPlayedAt := time.Unix(int64(log.Date), 0)

		if logPlayedAt.Before(matchPlayedAtMinusOffset) || logPlayedAt.After(matchPlayedAtPlusOffset) {
			logrus.Debugf("match log #%d didnt match based on time limits", log.Id)
			continue
		}

		if _, ok := mapsWhitelist[log.Map]; ok {
			matchLogs = append(matchLogs, log)
		} else {
			combinedLogs = append(combinedLogs, log)
		}
	}
	return matchLogs, combinedLogs
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
