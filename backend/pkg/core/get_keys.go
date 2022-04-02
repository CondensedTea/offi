package core

func (c Core) GetKeys(hashKey string) ([]string, error) {
	return c.cache.GetAllKeys(hashKey)
}
