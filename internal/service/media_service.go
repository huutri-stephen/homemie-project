package service

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

type MediaService interface {
	GeneratePresignedUploadURL(bucketName, objectKey string) (string, error)
	UploadFile(bucketName, objectKey string, file io.Reader, size int64) (string, error)
	CheckBucketName(bucketName string) error
}

type mediaService struct {
	s3Client       *s3.Client
	logger         *zap.Logger
	externalEndpoint string
}

func NewMediaService(s3Client *s3.Client, logger *zap.Logger, externalEndpoint string) MediaService {
	return &mediaService{
		s3Client:       s3Client,
		logger:         logger,
		externalEndpoint: externalEndpoint,
	}
}

func (s *mediaService) GeneratePresignedUploadURL(bucketName, objectKey string) (string, error) {
	presignClient := s3.NewPresignClient(s.s3Client)

	req, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = 15 * time.Minute
	})

	if err != nil {
		s.logger.Error("failed to generate presigned URL", zap.Error(err))
		return "", err
	}

	return req.URL, nil
}

func (s *mediaService) UploadFile(bucketName, objectKey string, file io.Reader, size int64) (string, error) {
	_, err := s.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(objectKey),
		Body:          file,
		ContentLength: aws.Int64(size),
		ACL:           "public-read", // Cho phép public access
	})
	if err != nil {
		s.logger.Error("failed to upload file", zap.Error(err))
		return "", err
	}

	url := fmt.Sprintf("%s/%s/%s", s.externalEndpoint, bucketName, objectKey)
	return url, nil
}

func (s *mediaService) CheckBucketName(bucketName string) error {
	_, err := s.s3Client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		s.logger.Error("Bucket "+ bucketName +" does not exist", zap.Error(err))
		return err
	}
	return nil
} 