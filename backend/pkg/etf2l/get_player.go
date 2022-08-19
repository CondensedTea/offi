package etf2l

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var ErrPlayerNotFound = errors.New(`player does not have an etf2l account`)

func (c Client) GetPlayer(id int) (Player, error) {
	url := fmt.Sprintf("https://api.etf2l.org/player/%d.json", id)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return Player{}, fmt.Errorf("failed to get player from etf2l api: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusInternalServerError {
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
