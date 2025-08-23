package entity

type User struct {
	Id         string  `json:"id"`
	Username   string  `json:"username"`
	ProfileImg *string `json:"profileImg"`

	Timed
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
