package entity

type Game struct {
	Id                     string `json:"id"`
	Title                  string `json:"title"`
	MaxPlayers             uint   `json:"maxPlayers"`
	CurrentRound           uint   `json:"currentRound"`
	MaxRounds              uint   `json:"maxRounds"`
	CurrentLeadingPlayerId string `json:"currentLeadingPlayerId"`
}

type GameRequest struct {
	Title      string `json:"title"`
	MaxPlayers uint   `json:"maxPlayers"`
	MaxRounds  uint   `json:"maxRounds"`
}
