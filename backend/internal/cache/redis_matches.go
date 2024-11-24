package cache

import "strconv"

func (r Redis) GetMatch(logId int) (MatchPage, error) {
	var mp MatchPage
	if err := r.client.HGet(logsKey, strconv.Itoa(logId)).Scan(&mp); err != nil {
		return MatchPage{}, err
	}
	return mp, nil
}

func (r Redis) SetMatch(logIds []int, matchPage *MatchPage) error {
	var err error

	for _, id := range logIds {
		if err = r.client.HSet(logsKey, strconv.Itoa(id), matchPage).Err(); err != nil {
			return err
		}
	}
	return nil
}
