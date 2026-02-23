package storage

import (
	"context"
	"io"
	"time"
)

type (
	ObjectInfo struct {
		Key          string    `json:"key"`
		Size         int64     `json:"size"`
		LastModified time.Time `json:"last_modified"`
	}
	ListObjectsOpts struct {
		Prefix    string
		Limit     int
		PageToken string
	}

	IStorageAdapter interface {
		PutObject(ctx context.Context, bucket, key string, body io.Reader, size int64, contentType string) error
		GeneratePresignedGet(ctx context.Context, bucket, key string, expire time.Duration) (string, error)
		ListObjects(ctx context.Context, bucket string, opts ListObjectsOpts) ([]ObjectInfo, string, error)
	}
)
