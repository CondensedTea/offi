package cache

type Cache interface {
	GetMatch(matchId int) (Match, error)
	SetMatch(matchId int, match *Match) error
	FlushMatch(matchId int) error

	GetPlayer(playerID string) (string, error)
	SetPlayer(playerID, steamID string) error
}
