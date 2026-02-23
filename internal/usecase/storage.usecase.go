package usecase

import (
	"context"
	"io"
	"path/filepath"
	"strings"
	"time"

	"go-structure/internal/adapter/storage"

	"github.com/google/uuid"
)

type ListInput struct {
	Bucket    string
	Prefix    string
	Limit     int
	PageToken string
}

const (
	uploadDirPrefix    = "upload"
	dateLayoutYYYYMMDD = "20060102"
)

type (
	UploadInput struct {
		Bucket           string
		Path             string
		Body             io.Reader
		Size             int64
		ContentType      string
		OriginalFilename string
	}

	UploadResult struct {
		Key string
	}

	ListResult struct {
		Items         []storage.ObjectInfo `json:"items"`
		NextPageToken string               `json:"next_page_token,omitempty"`
	}
)

type IStorageUsecase interface {
	PutObject(ctx context.Context, bucket, key string, body io.Reader, size int64, contentType string) error
	Upload(ctx context.Context, in UploadInput) (*UploadResult, error)
	List(ctx context.Context, in ListInput) (*ListResult, error)
	GeneratePresignedGet(ctx context.Context, bucket, key string, expire time.Duration) (string, error)
}

type storageUsecase struct {
	adapter storage.IStorageAdapter
}

func NewStorageUsecase(adapter storage.IStorageAdapter) IStorageUsecase {
	return &storageUsecase{adapter: adapter}
}

func (u *storageUsecase) PutObject(ctx context.Context, bucket, key string, body io.Reader, size int64, contentType string) error {
	return u.adapter.PutObject(ctx, bucket, key, body, size, contentType)
}

func (u *storageUsecase) Upload(ctx context.Context, in UploadInput) (*UploadResult, error) {
	path := strings.Trim(strings.TrimSpace(in.Path), "/")
	now := time.Now()
	dmy := now.Format(dateLayoutYYYYMMDD)
	ext := filepath.Ext(in.OriginalFilename)
	if ext == "" {
		ext = ".bin"
	}
	name := uuid.New().String() + ext

	var key string
	if path == "" {
		key = uploadDirPrefix + "/" + dmy + "/" + name
	} else {
		key = uploadDirPrefix + "/" + path + "/" + dmy + "/" + name
	}

	if err := u.adapter.PutObject(ctx, in.Bucket, key, in.Body, in.Size, in.ContentType); err != nil {
		return nil, err
	}
	return &UploadResult{Key: key}, nil
}

func (u *storageUsecase) List(ctx context.Context, in ListInput) (*ListResult, error) {
	items, nextToken, err := u.adapter.ListObjects(ctx, in.Bucket, storage.ListObjectsOpts{
		Prefix:    in.Prefix,
		Limit:     in.Limit,
		PageToken: in.PageToken,
	})
	if err != nil {
		return nil, err
	}
	return &ListResult{Items: items, NextPageToken: nextToken}, nil
}

func (u *storageUsecase) GeneratePresignedGet(ctx context.Context, bucket, key string, expire time.Duration) (string, error) {
	return u.adapter.GeneratePresignedGet(ctx, bucket, key, expire)
}
