package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"go-structure/global"
	"go-structure/internal/adapter/storage"
	"go-structure/internal/middleware"

	"github.com/google/uuid"
	"go.uber.org/zap"
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

	// MaxDriverDocumentSize là giới hạn kích thước file tài liệu tài xế (10 MB)
	MaxDriverDocumentSize int64 = 10 * 1024 * 1024
)

var (
	// allowedDriverDocumentExts là danh sách đuôi file hợp lệ cho tài liệu tài xế
	allowedDriverDocumentExts = map[string]struct{}{
		".jpg":  {},
		".jpeg": {},
		".png":  {},
		".webp": {},
		".heic": {},
		".pdf":  {},
		".doc":  {},
		".docx": {},
	}

	ErrDriverDocumentInvalidExtension = errors.New("định dạng file không được hỗ trợ, chỉ chấp nhận: jpg, jpeg, png, webp, heic, pdf, doc, docx")
	ErrDriverDocumentFileTooLarge     = fmt.Errorf("kích thước file vượt quá giới hạn cho phép (%dMB)", MaxDriverDocumentSize/1024/1024)
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
	UploadDriverDocument(ctx context.Context, in UploadInput, driverProfileID uuid.UUID) (*UploadResult, error)
}

type storageUsecase struct {
	adapter storage.IStorageAdapter
}

func NewStorageUsecase(adapter storage.IStorageAdapter) IStorageUsecase {
	return &storageUsecase{adapter: adapter}
}

func (u *storageUsecase) PutObject(ctx context.Context, bucket, key string, body io.Reader, size int64, contentType string) error {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("PutObject: start", zap.String(global.KeyCorrelationID, cid), zap.String("bucket", bucket), zap.String("key", key))
	if err := u.adapter.PutObject(ctx, bucket, key, body, size, contentType); err != nil {
		global.Logger.Error("PutObject: failed", zap.String(global.KeyCorrelationID, cid), zap.String("bucket", bucket), zap.String("key", key), zap.Error(err))
		return err
	}
	global.Logger.Info("PutObject: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("key", key))
	return nil
}

func (u *storageUsecase) Upload(ctx context.Context, in UploadInput) (*UploadResult, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("Upload: start", zap.String(global.KeyCorrelationID, cid), zap.String("bucket", in.Bucket), zap.String("original_filename", in.OriginalFilename))
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
		global.Logger.Error("Upload: failed to put object", zap.String(global.KeyCorrelationID, cid), zap.String("key", key), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("Upload: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("key", key))
	return &UploadResult{Key: key}, nil
}

func (u *storageUsecase) UploadDriverDocument(ctx context.Context, in UploadInput, driverProfileID uuid.UUID) (*UploadResult, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("UploadDriverDocument: start", zap.String(global.KeyCorrelationID, cid), zap.String("driver_profile_id", driverProfileID.String()), zap.String("original_filename", in.OriginalFilename))

	if in.Size > MaxDriverDocumentSize {
		global.Logger.Error("UploadDriverDocument: file too large", zap.String(global.KeyCorrelationID, cid), zap.Int64("size", in.Size), zap.Int64("max_size", MaxDriverDocumentSize))
		return nil, ErrDriverDocumentFileTooLarge
	}

	ext := strings.ToLower(filepath.Ext(in.OriginalFilename))
	if ext == "" {
		global.Logger.Error("UploadDriverDocument: missing file extension", zap.String(global.KeyCorrelationID, cid))
		return nil, ErrDriverDocumentInvalidExtension
	}
	if _, ok := allowedDriverDocumentExts[ext]; !ok {
		global.Logger.Error("UploadDriverDocument: invalid file extension", zap.String(global.KeyCorrelationID, cid), zap.String("ext", ext))
		return nil, ErrDriverDocumentInvalidExtension
	}

	now := time.Now()
	dmy := now.Format(dateLayoutYYYYMMDD)
	key := uploadDirPrefix + "/documents/" + dmy + "/drivers/" + driverProfileID.String() + "/" + uuid.New().String() + ext

	if err := u.adapter.PutObject(ctx, in.Bucket, key, in.Body, in.Size, in.ContentType); err != nil {
		global.Logger.Error("UploadDriverDocument: failed to put object", zap.String(global.KeyCorrelationID, cid), zap.String("key", key), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("UploadDriverDocument: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("key", key))
	return &UploadResult{Key: key}, nil
}

func (u *storageUsecase) List(ctx context.Context, in ListInput) (*ListResult, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("List: start", zap.String(global.KeyCorrelationID, cid), zap.String("bucket", in.Bucket), zap.String("prefix", in.Prefix))
	items, nextToken, err := u.adapter.ListObjects(ctx, in.Bucket, storage.ListObjectsOpts{
		Prefix:    in.Prefix,
		Limit:     in.Limit,
		PageToken: in.PageToken,
	})
	if err != nil {
		global.Logger.Error("List: failed to list objects", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("List: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.Int("count", len(items)))
	return &ListResult{Items: items, NextPageToken: nextToken}, nil
}

func (u *storageUsecase) GeneratePresignedGet(ctx context.Context, bucket, key string, expire time.Duration) (string, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("GeneratePresignedGet: start", zap.String(global.KeyCorrelationID, cid), zap.String("bucket", bucket), zap.String("key", key))
	url, err := u.adapter.GeneratePresignedGet(ctx, bucket, key, expire)
	if err != nil {
		global.Logger.Error("GeneratePresignedGet: failed", zap.String(global.KeyCorrelationID, cid), zap.String("key", key), zap.Error(err))
		return "", err
	}
	global.Logger.Info("GeneratePresignedGet: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("key", key))
	return url, nil
}
