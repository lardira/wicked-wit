package entity

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/lardira/wicked-wit/internal/db/model"
)

type Timed struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func TimedFromModel(model *model.Timed) Timed {
	return Timed{
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

type simpleData struct {
	Data any `json:"data"`
}

func SimpleData(w http.ResponseWriter, v any) {
	if err := json.NewEncoder(w).Encode(simpleData{Data: v}); err != nil {
		log.Printf("error when trying to respond with data: %v\n", v)
		SimpleError(w, err, http.StatusInternalServerError)
		return
	}
}

func SimpleError(w http.ResponseWriter, err error, code int) {
	if code < http.StatusContinue || code > http.StatusNetworkAuthenticationRequired {
		code = http.StatusInternalServerError
	}
	if err == nil {
		err = errors.New(http.StatusText(code))
	}
	http.Error(w, err.Error(), code)
}
