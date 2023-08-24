package web

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/likeawizard/document-ai-demo/config"
	"github.com/likeawizard/document-ai-demo/database"
	"github.com/likeawizard/document-ai-demo/expensebot"
	"github.com/likeawizard/document-ai-demo/store"
)

var Router *gin.Engine

// Define allowed file types. source: https://cloud.google.com/document-ai/docs/file-types
var supportedMimeTypes = []string{"application/pdf", "image/gif", "image/tiff", "image/jpeg", "image/png", "image/bmp", "image/webp"}

func NewRouter(cfg config.AppCfg) *gin.Engine {
	if !config.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.SetTrustedProxies(nil)

	expensesRest := router.Group("expenses")
	expensesRest.POST("", expensesCreate)
	expensesRest.GET(":uuid", expensesGetOne)

	return router
}

func expensesCreate(c *gin.Context) {
	id := uuid.New()
	formFile, _ := c.FormFile("file")
	mimeType := formFile.Header.Get("Content-Type")
	if !isSupportedMimeType(mimeType) {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("unsupported MIME Type '%s'", mimeType))
		return
	}
	f, err := formFile.Open()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer f.Close()

	newFilename := fmt.Sprintf("%s%s", id, filepath.Ext(formFile.Filename))
	store.File.Store(newFilename, f)

	record := database.New(id)
	record.Filename = formFile.Filename
	record.MimeType = mimeType
	record.Path = newFilename
	database.Instance.Create(record)
	err = expensebot.Processor.Process(record)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusOK, record)
}

func expensesGetOne(c *gin.Context) {
	id := c.Param("uuid")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	record, err := database.Instance.Get(uuid)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.IndentedJSON(http.StatusOK, record)
}

func isSupportedMimeType(mimeType string) bool {
	for _, supported := range supportedMimeTypes {
		if supported == mimeType {
			return true
		}
	}
	return false
}
