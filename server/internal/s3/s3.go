package s3

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	Client *S3Client
)

type S3Client struct {
	accessKey     string
	secretKey     string
	Url           string
	DefaultBucket string
	UseSsl        bool

	Conn *minio.Client
}

type S3Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Url             string
	Bucket          string
	UseSsl          bool
}

func Init(config *S3Config) error {
	ctx := context.Background()

	minioClient, err := minio.New(config.Url, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSsl,
	})
	if err != nil {
		return err
	}

	err = minioClient.MakeBucket(ctx, config.Bucket, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, config.Bucket)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s is already created\n", config.Bucket)
		} else {
			return err
		}
	} else {
		log.Printf("Successfully created %s\n", config.Bucket)
	}

	Client = &S3Client{
		accessKey:     config.AccessKeyID,
		secretKey:     config.SecretAccessKey,
		Url:           minioClient.EndpointURL().String(),
		DefaultBucket: config.Bucket,
		UseSsl:        config.UseSsl,
		Conn:          minioClient,
	}

	return nil
}
