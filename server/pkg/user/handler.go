package user

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lardira/wicked-wit/internal/s3"
	"github.com/lardira/wicked-wit/pkg/response"
	"github.com/minio/minio-go/v7"
)

const (
	maxFileSize = 32 << 20 // 32 MB
)

var (
	validImgContentTypes   = []string{"image/jpeg", "image/png", "image/webp"}
	validImgFileExtensions = []string{".jpg", ".jpeg", ".png", ".webp"}
)

type Handler struct{}

func Router() chi.Router {
	var handler Handler
	r := chi.NewRouter()

	r.Get("/{id}", handler.GetUser)
	r.Post("/", handler.CreateUser)
	r.Post("/{id}/image", handler.UpdateProfileImage)
	r.Delete("/{id}", handler.DeleteUser)

	return r
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := uuid.Validate(id); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, err := SelectUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	payload := User{
		Id:       user.Id,
		Username: user.Username,
		Timed:    response.TimedFromModel(&user.TimedModel),
	}

	if user.ProfileImg.Valid {
		payload.ProfileImg = &user.ProfileImg.String
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user UserRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	newId, err := InsertUser(
		user.Username,
		user.Password,
	)
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

func (h *Handler) UpdateProfileImage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		// TODO: check if user exists
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	r.ParseMultipartForm(maxFileSize)

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

	fileName := id + path.Ext(header.Filename)
	fileUrl, _ := url.JoinPath(s3.Client.Url, s3.Client.DefaultBucket, fileName)

	_, err = s3.Client.Conn.PutObject(
		context.Background(),
		os.Getenv("MINIO_BUCKET_NAME"),
		fileName,
		file,
		header.Size,
		minio.PutObjectOptions{ContentType: header.Header["Content-Type"][0]},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := UpdateUserImg(id, fileUrl); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SimpleData(w, fileUrl)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	DeleteUser(id)
	w.WriteHeader(http.StatusNoContent)
}
