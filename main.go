package main

import (
	"log"

	"github.com/likeawizard/document-ai-demo/config"
	"github.com/likeawizard/document-ai-demo/database"
	"github.com/likeawizard/document-ai-demo/expensebot"
	"github.com/likeawizard/document-ai-demo/store"
	"github.com/likeawizard/document-ai-demo/web"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v\n", err)
	}

	restService, err := web.NewRestService(cfg)
	if err != nil {
		log.Fatalf("failed to initialize restService: %s", err)
	}

	restService.Router.Run()
}

func initAll(cfg config.Config) {
	config.Init(cfg)

	fileStore, err := store.NewFileStore(cfg.Store)
	if err != nil {
		log.Fatalf("unable to initialize file store: %v\n", err)
	}

	store.File = fileStore
	db, err := database.NewDataBase(cfg.Db)
	if err != nil {
		log.Fatalf("failed to initialize database: %v\n", err)
	}

	database.Instance = db
	expensebot.Processor, err = expensebot.NewDocumentProcessor(cfg.Processor)
	if err != nil {
		log.Fatalf("failed to initialize document processor: %v\n", err)
	}
}
