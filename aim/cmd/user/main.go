package main

import (
	"fmt"
	"log"
	"net"

	"github.com/aim/aim/configs"
	"github.com/aim/aim/internal/pkg/database"
	"github.com/aim/aim/internal/user/handler"
	"github.com/aim/aim/internal/user/repo"
	"google.golang.org/grpc"
)

func main() {
	cfg := configs.Load()

	db, err := database.NewDB(cfg.Database)
	if err != nil {
		log.Fatalf("[user] db init error: %v", err)
	}
	if err := repo.NewUserRepo(db).AutoMigrate(); err != nil {
		log.Fatalf("[user] migrate error: %v", err)
	}

	h := handler.NewUserHandler(db)
	_ = h

	addr := fmt.Sprintf(":%d", cfg.Server.User.GRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("[user] listen error: %v", err)
	}

	srv := grpc.NewServer()
	log.Printf("[user] serving on %s", addr)
	log.Fatal(srv.Serve(lis))
}
