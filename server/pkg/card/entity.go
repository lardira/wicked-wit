package card

type PlayCardRequest struct {
	RoundId int    `json:"roundId"`
	CardIds []int  `json:"cardIds"`
	UserId  string `json:"userId"`
}
