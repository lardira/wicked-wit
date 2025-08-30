package entity

import (
	"github.com/lardira/wicked-wit/internal/domain/repository"
	"github.com/lardira/wicked-wit/internal/helper/response"
)

type Round struct {
	Id       int                    `json:"id"`
	WinnerId *string                `json:"winnerId"`
	Position int                    `json:"position"`
	GameId   string                 `json:"gameId"`
	Status   repository.RoundStatus `json:"status"`

	response.Timed
}
