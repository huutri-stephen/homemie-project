package infra

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewS3Client() *s3.Client {
	endpoint := os.Getenv("S3_ENDPOINT")
	accessKey := os.Getenv("S3_ACCESS_KEY_ID")
	secretKey := os.Getenv("S3_SECRET_ACCESS_KEY")
	region := os.Getenv("S3_REGION")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: endpoint, SigningRegion: region}, nil
			})),
	)
	if err != nil {
		log.Fatalf("failed to load S3 config: %v", err)
	}

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
}

func EnsureBucketPolicy(s3Client *s3.Client, bucketName string) error {
	policy := `{
	  "Version": "2012-10-17",
	  "Statement": [
	    {
	      "Action": [
	        "s3:GetBucketLocation",
	        "s3:ListBucket"
	      ],
	      "Effect": "Allow",
	      "Principal": "*",
	      "Resource": "arn:aws:s3:::` + bucketName + `"
	    },
	    {
	      "Action": "s3:GetObject",
	      "Effect": "Allow",
	      "Principal": "*",
	      "Resource": "arn:aws:s3:::` + bucketName + `/*"
	    }
	  ]
	}`

	_, err := s3Client.PutBucketPolicy(context.TODO(), &s3.PutBucketPolicyInput{
		Bucket: aws.String(bucketName),
		Policy: aws.String(policy),
	})
	return err
}

func CreateBucketIfNotExists(client *s3.Client, bucketName string) {
	ctx := context.TODO()

	_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err == nil {
		log.Printf("Bucket %s already exists.", bucketName)
		return
	}

	log.Printf("Bucket %s does not exist. Creating it...", bucketName)

	// Only add CreateBucketConfiguration if region != us-east-1
	_, createErr := client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if createErr != nil {
		log.Fatalf("failed to create bucket %s: %v", bucketName, createErr)
	}

	log.Printf("Bucket %s created successfully.", bucketName)
	EnsureBucketPolicy(client, bucketName)
	log.Printf("Public read policy applied to bucket %s", bucketName)
}

