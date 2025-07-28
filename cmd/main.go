package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"zip-archive/internal/config"
	"zip-archive/internal/service"
	proto "zip-archive/protos"
)

func main() {
	cfg, err := config.MustLoad("config/config.yaml")
	if err != nil {
		log.Fatalf("failed to read config: %v", err.Error())
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	s := grpc.NewServer()
	archiveZipService := service.New(cfg)
	proto.RegisterArchiveZipServiceServer(s, archiveZipService)

	log.Printf("gRPC server started on :%d", cfg.Server.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatal("Failed to serve:", err)
	}
}
