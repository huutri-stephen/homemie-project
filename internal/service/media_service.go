package service

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

type MediaService interface {
	GeneratePresignedUploadURL(bucketName, objectKey string) (string, error)
}

type mediaService struct {
	s3Client *s3.Client
	logger   *zap.Logger
}

func NewMediaService(s3Client *s3.Client, logger *zap.Logger) MediaService {
	return &mediaService{
		s3Client: s3Client,
		logger:   logger,
	}
}

func (s *mediaService) GeneratePresignedUploadURL(bucketName, objectKey string) (string, error) {
	presignClient := s3.NewPresignClient(s.s3Client)

	req, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(15 * time.Minute)
	})

	if err != nil {
		s.logger.Error("failed to generate presigned URL", zap.Error(err))
		return "", err
	}

	return req.URL, nil
}
