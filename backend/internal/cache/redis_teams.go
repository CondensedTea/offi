package cache

import (
	"fmt"
	"time"
)

const teamExpiration = 12 * time.Hour

func (r Redis) GetTeam(teamID int) (Team, error) {
	key := fmt.Sprintf("team-%d", teamID)

	var team Team
	if err := r.client.Get(key).Scan(&team); err != nil {
		return Team{}, err
	}
	return team, nil
}

func (r Redis) SetTeam(teamID int, team Team) error {
	key := fmt.Sprintf("team-%d", teamID)

	return r.client.Set(key, team, teamExpiration).Err()
}
