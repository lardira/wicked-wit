package entity

import (
	"github.com/lardira/wicked-wit/internal/domain/repository"
	"github.com/lardira/wicked-wit/internal/helper/response"
)

type Game struct {
	Id           string                `json:"id"`
	Title        string                `json:"title"`
	MaxPlayers   int                   `json:"maxPlayers"`
	MaxRound     int                   `json:"maxRound"`
	CurrentRound int                   `json:"currentRound"`
	Status       repository.GameStatus `json:"status"`
	UserHostId   string                `json:"userHostId"`

	response.Timed
}

type GameRequest struct {
	Title      string   `json:"title"`
	MaxPlayers int      `json:"maxPlayers"`
	MaxRound   int      `json:"maxRound"`
	Users      []string `json:"userIds"`
	HostId     string   `json:"hostId"`
}
