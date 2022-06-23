package core

import (
	"fmt"
	"offi/pkg/etf2l"
	"strconv"
)

func (c Core) GetPlayer(playerId string) (etf2l.Player, error) {
	id, err := strconv.Atoi(playerId)
	if err != nil {
		return etf2l.Player{}, fmt.Errorf("failed to parse player id (%s): %v", playerId, err)
	}

	return c.etf2l.GetPlayer(id)
}
