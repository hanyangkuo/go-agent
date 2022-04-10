package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	endpoint = "localhost:9000"
	accessKeyID = "admin"
	secretAccessKey = "password"
	useSSL = false
)

func Download(ctx context.Context, bucketName string, objectName string) error{
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return err
	}

	err = minioClient.FGetObject(context.Background(), bucketName, objectName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}