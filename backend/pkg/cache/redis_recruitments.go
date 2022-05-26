package cache

import (
	"fmt"
	"time"
)

const postExpiration = 2 * 24 * time.Hour

func (r Redis) SaveRecruitmentPosts(postType string, entries []Entry) error {
	for _, entry := range entries {
		key := fmt.Sprintf("recruiment-%s-%d", postType, entry.ID)
		if err := r.client.Set(key, entry, postExpiration).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (r Redis) GetRecruitmentPost(postType, id string) (*Entry, error) {
	key := fmt.Sprintf("recruiment-%s-%s", postType, id)

	var entry Entry
	if err := r.client.Get(key).Scan(&entry); err != nil {
		return nil, err
	}
	return &entry, nil
}
