package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/aim/aim/configs"
	"github.com/aim/aim/internal/ai/handler"
	airepo "github.com/aim/aim/internal/ai/repo"
	"github.com/aim/aim/internal/pkg/database"
	"google.golang.org/grpc"
)

func main() {
	cfg := configs.Load()

	db, err := database.NewDB(cfg.Database)
	if err != nil {
		log.Fatalf("[ai] db init error: %v", err)
	}
	if err := airepo.NewAIRepo(db).AutoMigrate(); err != nil {
		log.Fatalf("[ai] migrate error: %v", err)
	}

	h := handler.NewAIHandler(db, nil)
	_ = h

	grpcAddr := fmt.Sprintf(":%d", cfg.Server.AI.GRPCPort)
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("[ai] grpc listen error: %v", err)
	}
	grpcSrv := grpc.NewServer()
	go func() {
		log.Printf("[ai] gRPC serving on %s", grpcAddr)
		log.Fatal(grpcSrv.Serve(lis))
	}()

	httpAddr := fmt.Sprintf(":%d", cfg.Server.AI.HTTPPort)
	log.Printf("[ai] HTTP serving on %s", httpAddr)
	log.Fatal(http.ListenAndServe(httpAddr, nil))
}
