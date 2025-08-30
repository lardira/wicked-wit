package service

import (
	"context"
	"mime/multipart"
	"net/url"
	"os"
	"path"

	"github.com/lardira/wicked-wit/internal/domain/entity"
	"github.com/lardira/wicked-wit/internal/domain/repository"
	"github.com/lardira/wicked-wit/internal/helper/response"
	"github.com/lardira/wicked-wit/internal/s3"
	"github.com/minio/minio-go/v7"
)

type userService struct{}

func NewUserService() *userService {
	return &userService{}
}

func (s *userService) GetUser(id string) (*entity.User, error) {
	userModel, err := repository.SelectUser(id)
	if err != nil {
		return nil, err
	}

	user := entity.User{
		Id:       userModel.Id,
		Username: userModel.Username,
		Timed:    response.TimedFromModel(&userModel.TimedModel),
	}

	if userModel.ProfileImg.Valid {
		user.ProfileImg = &userModel.ProfileImg.String
	}

	return &user, err
}

func (s *userService) CreateUser(userRequest *entity.UserRequest) (string, error) {
	newId, err := repository.InsertUser(
		userRequest.Username,
		userRequest.Password,
	)
	if err != nil {
		return "", err
	}

	return newId, nil
}

func (s *userService) UpdateProfileImage(id string, file *multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	fileName := id + path.Ext(fileHeader.Filename)
	fileUrl, _ := url.JoinPath(s3.Client.Url, s3.Client.DefaultBucket, fileName)

	_, err := s3.Client.Conn.PutObject(
		context.Background(),
		os.Getenv("MINIO_BUCKET_NAME"),
		fileName,
		*file,
		fileHeader.Size,
		minio.PutObjectOptions{ContentType: fileHeader.Header["Content-Type"][0]},
	)
	if err != nil {
		return "", err
	}

	if err := repository.UpdateUserImg(id, fileUrl); err != nil {
		return "", err
	}

	return fileName, nil
}

func (s *userService) DeleteUser(id string) {
	repository.DeleteUser(id)
}
