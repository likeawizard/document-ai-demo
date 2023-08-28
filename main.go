package main

import (
	"log"

	"github.com/likeawizard/document-ai-demo/config"
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
