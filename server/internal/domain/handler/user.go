package handler

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"path"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lardira/wicked-wit/internal/domain/entity"
	"github.com/lardira/wicked-wit/internal/domain/interfaces"
	"github.com/lardira/wicked-wit/internal/helper/response"
)

const (
	maxFileSize = 32 << 20 // 32 MB
)

var (
	validImgContentTypes   = []string{"image/jpeg", "image/png", "image/webp"}
	validImgFileExtensions = []string{".jpg", ".jpeg", ".png", ".webp"}
)

type userHandler struct {
	userService interfaces.UserService
}

func NewHandler(userService interfaces.UserService) *userHandler {
	return &userHandler{
		userService: userService,
	}
}

func UserRouter(userService interfaces.UserService) chi.Router {
	handler := NewHandler(
		userService,
	)

	r := chi.NewRouter()

	r.Get("/{id}", handler.GetUser)
	r.Post("/", handler.CreateUser)
	r.Post("/{id}/image", handler.UpdateProfileImage)
	r.Delete("/{id}", handler.DeleteUser)

	return r
}

func (h *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := uuid.Validate(id); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	payload, err := h.userService.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userRequest entity.UserRequest
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	newId, err := h.userService.CreateUser(&userRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SimpleData(w, newId)
}

func validateUploadedFileHeader(header *multipart.FileHeader) error {
	contentType := header.Header["Content-Type"][0]
	fileName := header.Filename
	fileExtension := path.Ext(fileName)

	if !slices.Contains(validImgContentTypes, contentType) {
		return fmt.Errorf("%v is unsupported content type for image", contentType)
	}

	if !slices.Contains(validImgFileExtensions, fileExtension) {
		return fmt.Errorf("%v is unsupported file extension for image", fileExtension)
	}

	return nil
}

func (h *userHandler) UpdateProfileImage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		// TODO: check if user exists
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("imgFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := validateUploadedFileHeader(header); err != nil {
		response.SimpleError(w, err, http.StatusBadRequest)
		return
	}

	fileUrl, err := h.userService.UpdateProfileImage(
		id,
		&file,
		header,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SimpleData(w, fileUrl)
}

func (h *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	h.userService.DeleteUser(id)

	w.WriteHeader(http.StatusNoContent)
}
