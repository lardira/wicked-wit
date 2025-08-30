package entity

type Card struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
}

type TemplateCard struct {
	Id                int    `json:"id"`
	Text              string `json:"text"`
	PlaceholdersCount int    `json:"placeholdersCount"`
}

type PlayCardRequest struct {
	RoundId int    `json:"roundId"`
	CardIds []int  `json:"cardIds"`
	UserId  string `json:"userId"`
}
