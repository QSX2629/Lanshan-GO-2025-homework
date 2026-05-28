package main

import (
	"fmt"
	"log"
	"net"

	"github.com/aim/aim/configs"
	"github.com/aim/aim/internal/message/handler"
	"github.com/aim/aim/internal/message/repo"
	"github.com/aim/aim/internal/pkg/database"
	"google.golang.org/grpc"
)

func main() {
	cfg := configs.Load()

	db, err := database.NewDB(cfg.Database)
	if err != nil {
		log.Fatalf("[message] db init error: %v", err)
	}
	if err := repo.NewMessageRepo(db).AutoMigrate(); err != nil {
		log.Fatalf("[message] migrate error: %v", err)
	}

	h := handler.NewMessageHandler(db)
	_ = h

	addr := fmt.Sprintf(":%d", cfg.Server.Message.GRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("[message] listen error: %v", err)
	}

	srv := grpc.NewServer()
	log.Printf("[message] serving on %s", addr)
	log.Fatal(srv.Serve(lis))
}
