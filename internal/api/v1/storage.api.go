package v1

import (
	ctl "go-structure/internal/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

const API_MODULE_STORAGE = "storage"

type (
	StorageApi interface {
		InitStorageApi(router *gin.RouterGroup, storageController ctl.StorageController)
	}

	storageApi struct {
		storageController ctl.StorageController
	}
)

func NewStorageApi(storageController ctl.StorageController) StorageApi {
	return &storageApi{storageController: storageController}
}

func (s *storageApi) InitStorageApi(router *gin.RouterGroup, storageController ctl.StorageController) {
	g := router.Group(API_MODULE_STORAGE)
	{
		g.POST("upload", func(c *gin.Context) {
			resp := storageController.Upload(c)
			c.JSON(http.StatusOK, resp)
		})
		g.POST("upload/driver-document", func(c *gin.Context) {
			resp := storageController.UploadDriverDocument(c)
			c.JSON(http.StatusOK, resp)
		})
		g.GET("list", func(c *gin.Context) {
			resp := storageController.List(c)
			c.JSON(http.StatusOK, resp)
		})
		g.GET("view-private", func(c *gin.Context) {
			resp := storageController.ViewPrivate(c)
			c.JSON(http.StatusOK, resp)
		})
	}
}
