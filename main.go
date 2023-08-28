package main

import (
	"log"

	"github.com/likeawizard/document-ai-demo/config"
	"github.com/likeawizard/document-ai-demo/expense"
	"github.com/likeawizard/document-ai-demo/web"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v\n", err)
	}

	processorEngine, err := expense.NewExpenseEngine(cfg)
	if err != nil {
		log.Fatalf("failed to initialize processorEngine: %v\n", err)
	}
	go processorEngine.Listen()
	procChan := processorEngine.GetSendChan()

	restService, err := web.NewRestService(cfg, procChan)
	if err != nil {
		log.Fatalf("failed to initialize restService: %s", err)
	}

	restService.Router.Run()
}
