package cache

func (r Redis) GetPlayer(playerID string) (string, error) {
	return r.client.HGet(playerKey, playerID).Result()
}

func (r Redis) SetPlayer(playerID, steamID string) error {
	return r.client.HSet(playerKey, playerID, steamID).Err()
}
