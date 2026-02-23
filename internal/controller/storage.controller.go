package controller

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-structure/config"
	"go-structure/internal/common"
	"go-structure/internal/usecase"

	"github.com/gin-gonic/gin"
)

const (
	VisibilityPublic    = "public"
	VisibilityPrivate   = "private"
	QueryTypeAll        = "all"
	FormKeyFile         = "file"
	FormKeyFileBrackets = "file[]"
	MaxUploadFiles      = 20
)

type (
	StorageController interface {
		Upload(c *gin.Context) *common.ResponseData
		List(c *gin.Context) *common.ResponseData
		ViewPrivate(c *gin.Context) *common.ResponseData
	}

	storageController struct {
		*BaseController
		storageUsecase usecase.IStorageUsecase
		storageCfg     config.Storage
	}
)

func NewStorageController(storageUsecase usecase.IStorageUsecase, storageCfg config.Storage) StorageController {
	return &storageController{
		BaseController: NewBaseController(),
		storageUsecase: storageUsecase,
		storageCfg:     storageCfg,
	}
}

func (ctl *storageController) Upload(c *gin.Context) *common.ResponseData {
	form, err := c.MultipartForm()
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"multipart form required"})
	}

	files := form.File[FormKeyFile]
	if len(files) == 0 {
		files = form.File[FormKeyFileBrackets]
	}
	if len(files) == 0 && form.File != nil {
		for _, partFiles := range form.File {
			if len(partFiles) > 0 {
				files = partFiles
				break
			}
		}
	}
	if len(files) == 0 {
		return common.ErrorResponse(common.StatusBadRequest, []string{"file is required (form-data key: file hoặc file[])"})
	}
	if len(files) > MaxUploadFiles {
		return common.ErrorResponse(common.StatusBadRequest, []string{"too many files, max " + strconv.Itoa(MaxUploadFiles)})
	}

	path := strings.TrimSpace(c.PostForm("path"))
	visibility := strings.TrimSpace(strings.ToLower(c.PostForm("visibility")))
	if visibility == "" {
		visibility = VisibilityPublic
	}
	bucket := ctl.storageCfg.BucketPublic
	if visibility == VisibilityPrivate {
		bucket = ctl.storageCfg.BucketPrivate
	}

	uploads := make([]gin.H, 0, len(files))
	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			return common.ErrorResponse(http.StatusInternalServerError, []string{err.Error()})
		}
		contentType := file.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		result, err := ctl.storageUsecase.Upload(c.Request.Context(), usecase.UploadInput{
			Bucket:           bucket,
			Path:             path,
			Body:             f,
			Size:             file.Size,
			ContentType:      contentType,
			OriginalFilename: file.Filename,
		})
		_ = f.Close()
		if err != nil {
			return common.ErrorResponse(http.StatusInternalServerError, []string{err.Error()})
		}
		uploads = append(uploads, gin.H{
			"key":               result.Key,
			"path":              path,
			"size":              file.Size,
			"original_filename": file.Filename,
			"visibility":        visibility,
			"bucket":            bucket,
		})
	}

	return common.SuccessResponse(common.StatusOK, gin.H{"uploads": uploads})
}

func (ctl *storageController) List(c *gin.Context) *common.ResponseData {
	prefix := strings.TrimSpace(c.Query("prefix"))
	if strings.ToLower(c.Query("type")) == QueryTypeAll {
		prefix = ""
	}
	limit := 20
	if l := c.Query("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}
	pageToken := strings.TrimSpace(c.Query("page_token"))

	result, err := ctl.storageUsecase.List(c.Request.Context(), usecase.ListInput{
		Bucket:    ctl.storageCfg.BucketPublic,
		Prefix:    prefix,
		Limit:     limit,
		PageToken: pageToken,
	})
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, []string{err.Error()})
	}
	resp := gin.H{"items": result.Items}
	if result.NextPageToken != "" {
		resp["next_page_token"] = result.NextPageToken
	}
	return common.SuccessResponse(common.StatusOK, resp)
}

func (ctl *storageController) ViewPrivate(c *gin.Context) *common.ResponseData {
	key := c.Query("key")
	if key == "" {
		return common.ErrorResponse(common.StatusBadRequest, []string{"key required"})
	}

	// TODO: check quyền user ở đây

	expire := 15 * time.Minute

	url, err := ctl.storageUsecase.GeneratePresignedGet(
		c.Request.Context(),
		ctl.storageCfg.BucketPrivate,
		key,
		expire,
	)
	if err != nil {
		return common.ErrorResponse(http.StatusInternalServerError, []string{"cannot generate url"})
	}

	return common.SuccessResponse(common.StatusOK, gin.H{
		"url": url,
	})
}
