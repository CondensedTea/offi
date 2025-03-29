package cache

import (
	"context"
	"strconv"
)

func (r Redis) GetMatch(ctx context.Context, logId int) (MatchPage, error) {
	var mp MatchPage
	if err := r.client.HGet(ctx, logsKey, strconv.Itoa(logId)).Scan(&mp); err != nil {
		return MatchPage{}, err
	}
	return mp, nil
}

func (r Redis) SetMatch(ctx context.Context, logIds []int, matchPage *MatchPage) error {
	var err error

	for _, id := range logIds {
		if err = r.client.HSet(ctx, logsKey, strconv.Itoa(id), matchPage).Err(); err != nil {
			return err
		}
	}
	return nil
}
