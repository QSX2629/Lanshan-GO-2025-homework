package main

import (
	"fmt"
	"log"
	"net"

	"github.com/aim/aim/configs"
	"github.com/aim/aim/internal/pkg/database"
	"github.com/aim/aim/internal/relation/handler"
	"github.com/aim/aim/internal/relation/repo"
	"google.golang.org/grpc"
)

func main() {
	cfg := configs.Load()

	db, err := database.NewDB(cfg.Database)
	if err != nil {
		log.Fatalf("[relation] db init error: %v", err)
	}
	if err := repo.NewRelationRepo(db).AutoMigrate(); err != nil {
		log.Fatalf("[relation] migrate error: %v", err)
	}

	h := handler.NewRelationHandler(db)
	_ = h

	addr := fmt.Sprintf(":%d", cfg.Server.Relation.GRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("[relation] listen error: %v", err)
	}

	srv := grpc.NewServer()
	log.Printf("[relation] serving on %s", addr)
	log.Fatal(srv.Serve(lis))
}
