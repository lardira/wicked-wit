package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lardira/wicked-wit/entity"
	"github.com/lardira/wicked-wit/internal/db/model"
)

type UserHandler struct{}

func UserRouter() chi.Router {
	var handler UserHandler
	r := chi.NewRouter()

	r.Get("/{id}", handler.GetUser)
	r.Post("/", handler.CreateUser)
	r.Put("/{id}/image", handler.UpdateProfileImage)
	r.Delete("/{id}", handler.DeleteUser)

	return r
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := uuid.Validate(id); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, err := model.SelectUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	payload := entity.User{
		Id:       user.Id,
		Username: user.Username,
		Timed:    entity.TimedFromModel(&user.Timed),
	}

	if user.ProfileImg.Valid {
		payload.ProfileImg = &user.ProfileImg.String
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user entity.UserRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	newId, err := model.InsertUser(
		user.Username,
		user.Password,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	entity.SimpleData(w, newId)
}

// TODO: add s3
func (h *UserHandler) UpdateProfileImage(w http.ResponseWriter, r *http.Request) {
	entity.SimpleData(w, "TODO: add s3")
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	model.DeleteUser(id)
	w.WriteHeader(http.StatusNoContent)
}
