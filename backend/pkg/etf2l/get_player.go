package etf2l

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/carlmjohnson/requests"
)

var ErrPlayerNotFound = errors.New("etf2l api could not resolve player")

func GetPlayer(ctx context.Context, id int) (Player, error) {
	url := fmt.Sprintf("https://api.etf2l.org/player/%d.json", id)

	var playerResponse PlayerResponse
	err := requests.
		URL(url).
		ToJSON(&playerResponse).
		CheckStatus(http.StatusOK).
		Fetch(ctx)
	switch {
	case requests.HasStatusErr(err, http.StatusInternalServerError):
		return Player{}, ErrPlayerNotFound
	case err != nil:
		return Player{}, err
	}
	return playerResponse.Player, nil
}
