package web

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/likeawizard/document-ai-demo/config"
	"github.com/likeawizard/document-ai-demo/database"
	"github.com/likeawizard/document-ai-demo/expense"
	"github.com/likeawizard/document-ai-demo/store"
)

var Router *gin.Engine

type RestService struct {
	Router    *gin.Engine
	Db        database.DB
	EventChan expense.EventChan
	FileStore store.FileStore
}

// Define allowed file types. source: https://cloud.google.com/document-ai/docs/file-types
var supportedMimeTypes = []string{"application/pdf", "image/gif", "image/tiff", "image/jpeg", "image/png", "image/bmp", "image/webp"}

func NewRestService(cfg config.Config, eventChan expense.EventChan) (*RestService, error) {
	rest := RestService{
		Router:    NewRouter(cfg.App),
		EventChan: eventChan,
	}

	db, err := database.NewDataBase(cfg.Db)
	if err != nil {
		return nil, err
	}
	rest.Db = db

	store, err := store.NewFileStore(cfg.Store)
	if err != nil {
		return nil, err
	}
	rest.FileStore = store
	rest.registerRoutes()

	return &rest, nil

}

func (rest *RestService) registerRoutes() {
	expenses := rest.Router.Group("expenses")
	expenses.POST("", rest.expensesCreate)
	expenses.GET(":uuid", rest.expensesGetOne)
}

func NewRouter(cfg config.AppCfg) *gin.Engine {
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.SetTrustedProxies(nil)

	return router
}

func (rest *RestService) expensesCreate(c *gin.Context) {
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
	rest.FileStore.Store(newFilename, f)

	record := database.New(id)
	record.Filename = formFile.Filename
	record.MimeType = mimeType
	record.Path = newFilename
	rest.Db.Create(record)

	rest.EventChan.MsgNew(record)

	c.IndentedJSON(http.StatusOK, record)
}

func (rest *RestService) expensesGetOne(c *gin.Context) {
	id := c.Param("uuid")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	record, err := rest.Db.Get(uuid)
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
