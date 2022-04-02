package cache

type Cache interface {
	GetLogs(matchId int) (LogSet, error)
	SetLogs(matchId int, match *LogSet) error

	GetPlayer(playerID string) (string, error)
	SetPlayer(playerID, steamID string) error

	GetMatch(logId int) (MatchPage, error)
	SetMatch(logIds []int, matchPage *MatchPage) error

	GetAllKeys(hashKey string) ([]string, error)
}
