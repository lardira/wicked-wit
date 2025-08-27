package user

import (
	"context"
	"mime/multipart"
	"net/url"
	"os"
	"path"

	"github.com/lardira/wicked-wit/internal/s3"
	"github.com/lardira/wicked-wit/pkg/response"
	"github.com/minio/minio-go/v7"
)

type Service struct{}

func (s *Service) GetUser(id string) (*User, error) {
	userModel, err := SelectUser(id)
	if err != nil {
		return nil, err
	}

	user := User{
		Id:       userModel.Id,
		Username: userModel.Username,
		Timed:    response.TimedFromModel(&userModel.TimedModel),
	}

	if userModel.ProfileImg.Valid {
		user.ProfileImg = &userModel.ProfileImg.String
	}

	return &user, err
}

func (s *Service) CreateUser(userRequest *UserRequest) (string, error) {
	newId, err := InsertUser(
		userRequest.Username,
		userRequest.Password,
	)
	if err != nil {
		return "", err
	}

	return newId, nil
}

func (s *Service) UpdateProfileImage(id string, file *multipart.File, fileHeader *multipart.FileHeader) (string, error) {
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

	if err := UpdateUserImg(id, fileUrl); err != nil {
		return "", err
	}

	return fileName, nil
}

func (s *Service) DeleteUser(id string) {
	DeleteUser(id)
}
