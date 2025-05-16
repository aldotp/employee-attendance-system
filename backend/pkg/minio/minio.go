package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageInterface interface {
	Upload(ctx context.Context, object *FileUploadObject) error
	Download(ctx context.Context, fileName string) (*bytes.Buffer, error)
	Delete(ctx context.Context, fileName string) error
	GenerateUrl(fileName string) string
	GeneratePresignedUrl(fileName string, expiration time.Duration) (string, error)
}

// FileUploadObject represents a file to be uploaded
type FileUploadObject struct {
	File     io.Reader
	FileName string
}

// MinioClient implements StorageInterface for MinIO
type MinioClient struct {
	client *minio.Client
	bucket string
}

// NewMinioClient creates a new MinIO client
func NewMinioClient(endpoint, accessKey, secretKey, bucketName string, useSSL bool) (*MinioClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	return &MinioClient{
		client: client,
		bucket: bucketName,
	}, nil
}

// Upload uploads a file to MinIO
func (m *MinioClient) Upload(ctx context.Context, object *FileUploadObject) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	_, err := m.client.PutObject(ctx, m.bucket, object.FileName, object.File, -1, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to upload object: %w", err)
	}

	return nil
}

// Download downloads a file from MinIO
func (m *MinioClient) Download(ctx context.Context, fileName string) (*bytes.Buffer, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	object, err := m.client.GetObject(ctx, m.bucket, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	defer object.Close()

	buff := new(bytes.Buffer)
	_, err = io.Copy(buff, object)
	if err != nil {
		return nil, fmt.Errorf("failed to copy object to buffer: %w", err)
	}

	return buff, nil
}

// Delete deletes a file from MinIO
func (m *MinioClient) Delete(ctx context.Context, fileName string) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	err := m.client.RemoveObject(ctx, m.bucket, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

func (m *MinioClient) GenerateUrl(fileName string) string {
	return fmt.Sprintf("%s/%s/%s", m.client.EndpointURL(), m.bucket, fileName)
}

func (m *MinioClient) GeneratePresignedUrl(fileName string, expiration time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	url, err := m.client.PresignedPutObject(ctx, m.bucket, fileName, expiration)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}

var _ StorageInterface = &MinioClient{}
