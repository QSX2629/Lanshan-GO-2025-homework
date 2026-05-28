package es

import (
	"context"
	"log"

	"github.com/elastic/go-elasticsearch/v8"

	"AIM/internal/common/logger"
)

var Client *elasticsearch.Client
var Ctx = context.Background()

func Init() {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://127.0.0.1:9200"},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal("ES connect failed:", err)
	}

	Client = es
	logger.Info("ES connected ✅")
}
