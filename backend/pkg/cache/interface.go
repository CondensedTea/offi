package cache

type Cache interface {
	GetLogs(matchId int) (LogSet, error)
	SetLogs(matchId int, match *LogSet) error

	SetLogError(matchId int, err error) error
	CheckLogError(matchId int) error

	DeleteLogs(matchId int) (*LogSet, error)

	GetPlayer(playerID int) (Player, error)
	SetPlayer(playerID int, player Player) error

	GetTeam(teamID int) (Team, error)
	SetTeam(teamID int, team Team) error

	GetMatch(logId int) (MatchPage, error)
	SetMatch(logIds []int, matchPage *MatchPage) error

	GetAllKeys(hashKey string) ([]string, error)

	IncrementViews(object string, id int) (int64, error)
	GetViews(object string, id int) (int64, error)
}
