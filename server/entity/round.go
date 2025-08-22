package entity

type Round struct {
	Id       int     `json:"id"`
	WinnerId *string `json:"winnerId"`
	Position int     `json:"position"`
	GameId   string  `json:"gameId"`

	Timed
}
