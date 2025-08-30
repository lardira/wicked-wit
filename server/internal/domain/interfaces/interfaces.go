package interfaces

import (
	"mime/multipart"

	"github.com/lardira/wicked-wit/internal/domain/entity"
)

type GameService interface {
	GetGames() ([]entity.Game, error)
	CreateGame(*entity.GameRequest) (string, error)
	DeleteGame(string)
}

type RoundService interface {
	GetRounds(string) ([]entity.Round, error)
	AddRound(string) (int, error)
	DeleteRound(int)
}

type CardService interface {
	GetCards(gameId string, userId string) ([]entity.Card, error)
	UseCards(gameId string, userId string, cardIds ...int) error
	PlayCards(roundId int, userId string, cardIds ...int) (int, error)
	GetUnusedTemplateCards(gameId string) ([]entity.TemplateCard, error)
	GetRandomTemplateCard(gameId string) (*entity.TemplateCard, error)
}

type UserService interface {
	GetUser(id string) (*entity.User, error)
	CreateUser(userRequest *entity.UserRequest) (string, error)
	UpdateProfileImage(id string, file *multipart.File, fileHeader *multipart.FileHeader) (string, error)
	DeleteUser(id string)
}
