package etf2l

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var ErrPlayerNotFound = errors.New(`player does not have an etf2l account`)

func (c Client) GetPlayer(ctx context.Context, steamID int64) (Player, error) {
	ctx, span := c.tracer.Start(ctx, "etf2l.GetPlayer",
		trace.WithAttributes(attribute.Int64("steam_id", steamID)),
	)
	defer span.End()

	t := time.Now()

	if err := c.limiter.Wait(ctx); err != nil {
		return Player{}, err
	}

	span.SetAttributes(attribute.Float64("rate_limit_wait", float64(time.Since(t).Milliseconds())))

	url := fmt.Sprintf("%s/player/%d", c.apiURL, steamID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return Player{}, fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Player{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return Player{}, ErrPlayerNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return Player{}, fmt.Errorf("non-200 status: %d", resp.StatusCode)
	}

	var playerResponse PlayerResponse
	if err = json.NewDecoder(resp.Body).Decode(&playerResponse); err != nil {
		return Player{}, err
	}

	return playerResponse.Player, nil
}
