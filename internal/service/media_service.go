package service

import (
	"context"
	"fmt"
	"io"
	"homemie/pkg/logger"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type IMediaService interface {
	GeneratePresignedUploadURL(ctx context.Context, bucketName, objectKey string) (string, error)
	UploadFile(ctx context.Context, bucketName, objectKey string, file io.Reader, size int64) (string, error)
	CheckBucketName(ctx context.Context, bucketName string) error
}

type mediaService struct {
	s3Client         *s3.Client
	externalEndpoint string
}

func NewMediaService(s3Client *s3.Client, externalEndpoint string) IMediaService {
	return &mediaService{
		s3Client: s3Client, externalEndpoint: externalEndpoint}
}

func (s *mediaService) GeneratePresignedUploadURL(ctx context.Context, bucketName, objectKey string) (string, error) {
	presignClient := s3.NewPresignClient(s.s3Client)

	req, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = 15 * time.Minute
	})

	if err != nil {
		logger.Errorw(ctx, "failed to generate presigned URL", "error", err)
		return "", err
	}

	return req.URL, nil
}

func (s *mediaService) UploadFile(ctx context.Context, bucketName, objectKey string, file io.Reader, size int64) (string, error) {
	_, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(objectKey),
		Body:          file,
		ContentLength: aws.Int64(size),
		ACL:           "public-read", // Cho phép public access
	})
	if err != nil {
		logger.Errorw(ctx, "failed to upload file", "error", err)
		return "", err
	}

	url := fmt.Sprintf("%s/%s/%s", s.externalEndpoint, bucketName, objectKey)
	return url, nil
}

func (s *mediaService) CheckBucketName(ctx context.Context, bucketName string) error {
	_, err := s.s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		logger.Errorw(ctx, "Bucket "+bucketName+" does not exist", "error", err)
		return err
	}
	return nil
}