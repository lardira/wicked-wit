package user

import "github.com/lardira/wicked-wit/pkg/response"

type User struct {
	Id         string  `json:"id"`
	Username   string  `json:"username"`
	ProfileImg *string `json:"profileImg"`

	response.Timed
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
