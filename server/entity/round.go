package entity

type RoundStatus uint

const (
	RoundStatusStarted            = 0
	RoundStatusVoting             = 1
	RoundStatusWinnerPresentation = 2
	RoundStatusEnded              = 3
)

type Round struct {
	Id       int         `json:"id"`
	WinnerId *string     `json:"winnerId"`
	Position int         `json:"position"`
	GameId   string      `json:"gameId"`
	Status   RoundStatus `json:"status"`

	Timed
}
