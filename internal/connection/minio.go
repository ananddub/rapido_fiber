package connection

import (
	"context"
	"fmt"

	"encore.app/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewS3Client(cfg *config.Config) (*s3.Client, error) {
	// Create a custom resolver for MinIO
	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           cfg.MinioEndpoint,
			SigningRegion: cfg.MinioRegion,
		}, nil
	})

	sdkConfig, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(cfg.MinioRegion),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.MinioUser, cfg.MinioPassword, "")),
		awsconfig.WithEndpointResolverWithOptions(resolver),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	// Create S3 client
	// UsePathStyle is deprecated in v2 and usually handled by the resolver/custom endpoint logic, but for MinIO specifically we ensure address correctness.
	// In v2, we just create the client from config.
	client := s3.NewFromConfig(sdkConfig, func(o *s3.Options) {
		o.UsePathStyle = true // Required for MinIO
	})

	return client, nil
}
