package storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gocloud.dev/blob"
	"gocloud.dev/blob/s3blob"
)

func NewMinioBlobStorage(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*blob.Bucket, error) {
	// Create a custom resolver that always resolves to the provided endpoint
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           endpoint,
			SigningRegion: "us-east-1",
		}, nil
	})

	// Configure the S3 client
	cfg := aws.Config{
		Credentials:                 credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		EndpointResolverWithOptions: customResolver,
		Region:                      "us-east-1",
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true // Required for Minio
	})

	// Initialize minio client object.
	blob, err := s3blob.OpenBucketV2(context.Background(), client, bucket, nil)
	if err != nil {
		return nil, err
	}

	return blob, nil
}
