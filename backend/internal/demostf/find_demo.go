package demostf

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var ErrNotFound = errors.New("demo not found")

type FindDemoRequest struct {
	PlayerIDs []int
	Map       string
	PlayedAt  time.Time
}
type Demo struct {
	ID int
}

func (c *Client) FindDemo(ctx context.Context, req FindDemoRequest) (Demo, error) {
	ctx, span := c.tracer.Start(ctx, "demostf.FindDemo")
	defer span.End()

	reqURL := url.URL{
		Scheme: "https",
		Host:   "api.demos.tf",
		Path:   "/demos/",
	}

	var playerIDsBuilder strings.Builder
	for _, playerID := range req.PlayerIDs {
		if playerIDsBuilder.Len() > 0 {
			playerIDsBuilder.WriteRune(',')
		}
		playerIDsBuilder.WriteString(strconv.Itoa(playerID))
	}

	query := url.Values{}
	query.Set("players", playerIDsBuilder.String())
	query.Set("after", strconv.FormatInt(req.PlayedAt.Unix(), 10))
	query.Set("before", strconv.FormatInt(req.PlayedAt.Add(2*time.Minute).Unix(), 10)) // window of 2 minutes for uploading demo

	reqURL.RawQuery = query.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), http.NoBody)
	if err != nil {
		return Demo{}, fmt.Errorf("building request: %w", err)
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return Demo{}, fmt.Errorf("doing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Demo{}, fmt.Errorf("api returned non-200 status code: %d", resp.StatusCode)
	}

	var r []Demo
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return Demo{}, fmt.Errorf("decoding reponse: %w", err)
	}

	if len(r) == 0 {
		return Demo{}, ErrNotFound
	}

	if len(r) > 1 {
		return Demo{}, fmt.Errorf("unexecpted number of demos found for log: %v", r)
	}

	return r[0], nil
}
