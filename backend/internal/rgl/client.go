package rgl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Client struct {
	client *http.Client
	tracer trace.Tracer
}

func NewClient(rt http.RoundTripper) *Client {
	return &Client{
		client: &http.Client{Transport: rt},
		tracer: otel.Tracer("rgl"),
	}
}

type Player struct {
	SteamID int64  `json:"steamId,string"`
	Name    string `json:"name"`
}

func (c *Client) GetPlayers(ctx context.Context, playerIDs []int64) ([]Player, error) {
	ctx, span := c.tracer.Start(ctx, "rgl.GetPlayers")
	defer span.End()

	stringPlayerIDs := make([]string, len(playerIDs))
	for i, id := range playerIDs {
		stringPlayerIDs[i] = strconv.FormatInt(id, 10)
	}

	reqBytes, err := json.Marshal(stringPlayerIDs)
	if err != nil {
		return nil, fmt.Errorf("marshalling player IDs: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.rgl.gg/v0/profile/getmany", bytes.NewReader(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return []Player{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var players []Player
	if err = json.NewDecoder(resp.Body).Decode(&players); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return players, nil
}
