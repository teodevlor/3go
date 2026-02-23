package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const defaultRegion = "us-east-1"

type S3Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Provider  string
}

type s3Adapter struct {
	client        *s3.Client
	presignClient *s3.PresignClient
}

func NewS3Adapter(cfg S3Config) (IStorageAdapter, error) {
	awsCfg := aws.Config{
		Region: defaultRegion,
		Credentials: credentials.NewStaticCredentialsProvider(
			cfg.AccessKey,
			cfg.SecretKey,
			"",
		),
	}

	var clientOpts []func(*s3.Options)
	if cfg.Endpoint != "" {
		awsCfg.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:               cfg.Endpoint,
					SigningRegion:     defaultRegion,
					HostnameImmutable: true,
				}, nil
			},
		)
		clientOpts = append(clientOpts, func(o *s3.Options) {
			o.UsePathStyle = true
		})
	}

	client := s3.NewFromConfig(awsCfg, clientOpts...)
	presignClient := s3.NewPresignClient(client)

	return &s3Adapter{
		client:        client,
		presignClient: presignClient,
	}, nil
}

func (a *s3Adapter) PutObject(ctx context.Context, bucket, key string, body io.Reader, size int64, contentType string) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   body,
	}
	if size >= 0 {
		input.ContentLength = aws.Int64(size)
	}
	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}
	_, err := a.client.PutObject(ctx, input)
	if err != nil {
		return fmt.Errorf("put object: %w", err)
	}
	return nil
}

func (a *s3Adapter) GeneratePresignedGet(ctx context.Context, bucket, key string, expire time.Duration) (string, error) {
	req, err := a.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expire))
	if err != nil {
		return "", fmt.Errorf("presign get object: %w", err)
	}
	return req.URL, nil
}

func (a *s3Adapter) ListObjects(ctx context.Context, bucket string, opts ListObjectsOpts) ([]ObjectInfo, string, error) {
	limit := opts.Limit
	if limit <= 0 {
		limit = 100
	}
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int32(int32(limit)),
	}
	if opts.Prefix != "" {
		input.Prefix = aws.String(opts.Prefix)
	}
	if opts.PageToken != "" {
		input.ContinuationToken = aws.String(opts.PageToken)
	}
	page, err := a.client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, "", fmt.Errorf("list objects: %w", err)
	}
	var out []ObjectInfo
	for _, obj := range page.Contents {
		var lastMod time.Time
		if obj.LastModified != nil {
			lastMod = *obj.LastModified
		}
		if obj.Key == nil || (len(*obj.Key) > 0 && (*obj.Key)[len(*obj.Key)-1] == '/') {
			continue
		}
		out = append(out, ObjectInfo{
			Key:          aws.ToString(obj.Key),
			Size:         int64(aws.ToInt64(obj.Size)),
			LastModified: lastMod,
		})
	}
	nextToken := aws.ToString(page.NextContinuationToken)
	return out, nextToken, nil
}
