package entity

type CardType string

const (
	CardTypeAnswer   CardType = "answer"
	CardTypeTemplate CardType = "template"
)

const (
	CardPlaceholder = '_'

	CardPlaceholdersMinCount = 1
	CardPlaceholdersMaxCount = 3
)

type CardTemplate struct {
	Id                int    `json:"id"`
	PlaceholdersCount int    `json:"placeholdersCount"`
	Text              string `json:"text"`
}

type CardAnswer struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
}

type CardRequest struct {
	Text string   `json:"text"`
	Type CardType `json:"type"`
}

type CardsUsedRequest struct {
	GameId string `json:"gameId"`
	Cards  []int  `json:"cards"`
	UserId string `json:"userId"`
}
