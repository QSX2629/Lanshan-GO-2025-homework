package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/aim/aim/configs"
	filerepo "github.com/aim/aim/internal/file/repo"
	"github.com/aim/aim/internal/pkg/database"
	"google.golang.org/grpc"
)

func main() {
	cfg := configs.Load()

	db, err := database.NewDB(cfg.Database)
	if err != nil {
		log.Fatalf("[file] db init error: %v", err)
	}
	if err := filerepo.NewFileRepo(db).AutoMigrate(); err != nil {
		log.Fatalf("[file] migrate error: %v", err)
	}

	grpcAddr := fmt.Sprintf(":%d", cfg.Server.File.GRPCPort)
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("[file] grpc listen error: %v", err)
	}
	grpcSrv := grpc.NewServer()
	go func() {
		log.Printf("[file] gRPC serving on %s", grpcAddr)
		log.Fatal(grpcSrv.Serve(lis))
	}()

	httpAddr := fmt.Sprintf(":%d", cfg.Server.File.HTTPPort)
	log.Printf("[file] HTTP serving on %s", httpAddr)
	log.Fatal(http.ListenAndServe(httpAddr, nil))
}
