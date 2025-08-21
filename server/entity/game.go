package entity

type Game struct {
	Id           string `json:"id"`
	Title        string `json:"title"`
	MaxPlayers   uint   `json:"maxPlayers"`
	CurrentRound *uint  `json:"currentRound"`
	MaxRound     uint   `json:"maxRound"`
}

type GameRequest struct {
	Title      string `json:"title"`
	MaxPlayers uint   `json:"maxPlayers"`
	MaxRound   uint   `json:"maxRound"`
}
