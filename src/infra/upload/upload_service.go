package upload

import (
	"TTCS/src/common/configs"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"mime/multipart"
	"time"
)

type UploadService struct {
	minioClient *minio.Client
	bucketName  string
}

func NewUploadService(minioClient *minio.Client) *UploadService {
	return &UploadService{
		minioClient: minioClient,
		bucketName:  configs.GetConfig().Minio.Bucket,
	}
}

func (u *UploadService) UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error) {
	buffer, errFile := file.Open()
	if errFile != nil {
		return "", errFile
	}

	defer func(buffer multipart.File) {
		_ = buffer.Close()
		if err := recover(); err != nil {

		}
	}(buffer)

	objectName := fmt.Sprintf("%s-%s", file.Filename, time.Now().Format("20060102150405"))
	contentType := file.Header["Content-Type"][0]

	_, errPut := u.minioClient.PutObject(ctx, configs.GetConfig().Minio.Bucket, objectName, buffer, file.Size, minio.PutObjectOptions{ContentType: contentType})
	if errPut != nil {
		return "", errPut
	}
	fileURL := fmt.Sprintf("http://%s/%s/%s", u.minioClient.EndpointURL().Host, configs.GetConfig().Minio.Bucket, objectName)
	return fileURL, nil
}
