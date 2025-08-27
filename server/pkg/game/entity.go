package game

import "github.com/lardira/wicked-wit/pkg/response"

type GameStatus uint

const (
	GameStatusStarted    = 0
	GameStatusInProgress = 1
	GameStatusEnded      = 2
)

type Game struct {
	Id         string     `json:"id"`
	Title      string     `json:"title"`
	MaxPlayers uint       `json:"maxPlayers"`
	MaxRound   uint       `json:"maxRound"`
	Status     GameStatus `json:"status"`
	UserHostId string     `json:"userHostId"`

	response.Timed
}

type GameRequest struct {
	Title      string `json:"title"`
	MaxPlayers uint   `json:"maxPlayers"`
	MaxRound   uint   `json:"maxRound"`
}
