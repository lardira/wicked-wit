package card

type Card struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
}

type PlayCardRequest struct {
	RoundId int    `json:"roundId"`
	CardIds []int  `json:"cardIds"`
	UserId  string `json:"userId"`
}
