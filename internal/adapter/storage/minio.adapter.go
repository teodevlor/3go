package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type (
	MinIOConfig struct {
		Endpoint      string
		AccessKey     string
		SecretKey     string
		BucketPublic  string
		BucketPrivate string
	}

	minioAdapter struct {
		client *minio.Client
	}
)

const (
	minioInitTimeout = 5 * time.Second
)

func NewMinIOAdapter(cfg MinIOConfig) (IStorageAdapter, error) {
	endpoint, useSSL, err := parseEndpoint(cfg.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse endpoint: %w", err)
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio new client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), minioInitTimeout)
	defer cancel()
	for _, bucket := range []string{cfg.BucketPublic, cfg.BucketPrivate} {
		if bucket == "" {
			continue
		}
		exists, err := client.BucketExists(ctx, bucket)
		if err != nil {
			return nil, fmt.Errorf("minio check bucket %q: %w", bucket, err)
		}
		if !exists {
			if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
				return nil, fmt.Errorf("minio create bucket %q: %w", bucket, err)
			}
		}
	}

	return &minioAdapter{client: client}, nil
}

func parseEndpoint(raw string) (host string, useSSL bool, err error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "localhost:9000", false, nil
	}
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		raw = "http://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", false, err
	}
	host = u.Host
	if u.Scheme == "https" {
		useSSL = true
	}
	return host, useSSL, nil
}

func (a *minioAdapter) PutObject(ctx context.Context, bucket, key string, body io.Reader, size int64, contentType string) error {
	opts := minio.PutObjectOptions{}
	if contentType != "" {
		opts.ContentType = contentType
	}
	_, err := a.client.PutObject(ctx, bucket, key, body, size, opts)
	if err != nil {
		return fmt.Errorf("put object: %w", err)
	}
	return nil
}

func (a *minioAdapter) GeneratePresignedGet(ctx context.Context, bucket, key string, expire time.Duration) (string, error) {
	presignedURL, err := a.client.PresignedGetObject(ctx, bucket, key, expire, nil)
	if err != nil {
		return "", fmt.Errorf("presigned get object: %w", err)
	}
	return presignedURL.String(), nil
}

func (a *minioAdapter) ListObjects(ctx context.Context, bucket string, opts ListObjectsOpts) ([]ObjectInfo, string, error) {
	limit := opts.Limit
	if limit <= 0 {
		limit = 100
	}
	listOpts := minio.ListObjectsOptions{Prefix: opts.Prefix, Recursive: true}
	if opts.PageToken != "" {
		listOpts.StartAfter = opts.PageToken
	}
	var out []ObjectInfo
	for obj := range a.client.ListObjects(ctx, bucket, listOpts) {
		if obj.Err != nil {
			return nil, "", fmt.Errorf("list objects: %w", obj.Err)
		}
		out = append(out, ObjectInfo{
			Key:          obj.Key,
			Size:         obj.Size,
			LastModified: obj.LastModified,
		})
		if len(out) >= limit {
			break
		}
	}
	var nextToken string
	if len(out) == limit {
		nextToken = out[len(out)-1].Key
	}
	return out, nextToken, nil
}
