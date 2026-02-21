package storage

import (
	"context"
	"fmt"
	"time"

	"tms-core-service/internal/domain/service"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3Storage struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucket        string
	expiry        time.Duration
}

// NewS3StorageService creates a new S3 storage service
func NewS3StorageService(region, bucket, accessKey, secretKey string, expiry time.Duration) service.StorageService {
	cfg := aws.Config{
		Region:      region,
		Credentials: credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
	}
	client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(client)

	return &s3Storage{
		client:        client,
		presignClient: presignClient,
		bucket:        bucket,
		expiry:        expiry,
	}
}

// GenerateUploadURL creates a presigned PUT URL for uploading a file
func (s *s3Storage) GenerateUploadURL(ctx context.Context, key string, contentType string) (string, error) {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}

	presigned, err := s.presignClient.PresignPutObject(ctx, input, s3.WithPresignExpires(s.expiry))
	if err != nil {
		return "", fmt.Errorf("s3: presign put object: %w", err)
	}

	return presigned.URL, nil
}

// GenerateDownloadURL creates a presigned GET URL for downloading a file
func (s *s3Storage) GenerateDownloadURL(ctx context.Context, key string) (string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	presigned, err := s.presignClient.PresignGetObject(ctx, input, s3.WithPresignExpires(s.expiry))
	if err != nil {
		return "", fmt.Errorf("s3: presign get object: %w", err)
	}

	return presigned.URL, nil
}
