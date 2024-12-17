package etf2l

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var ErrPlayerNotFound = errors.New(`player does not have an etf2l account`)

func (c Client) GetPlayer(ctx context.Context, id int) (Player, error) {
	if err := c.limiter.Wait(ctx); err != nil {
		return Player{}, err
	}

	url := fmt.Sprintf("%s/player/%d", c.apiURL, id)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return Player{}, fmt.Errorf("failed to get player from etf2l api: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return Player{}, ErrPlayerNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return Player{}, fmt.Errorf("etf2l api returned non-200 status: %d", resp.StatusCode)
	}

	var playerResponse PlayerResponse
	if err = json.NewDecoder(resp.Body).Decode(&playerResponse); err != nil {
		return Player{}, err
	}

	return playerResponse.Player, nil
}
