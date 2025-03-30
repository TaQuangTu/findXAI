package main

import (
	"findx/config"
	"findx/internal/server"
	"findx/pkg/protogen"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.NewConfig()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	searchServer := server.NewSearchServer(cfg.POSTGRES_DSN)
	protogen.RegisterSearchServiceServer(s, searchServer)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
