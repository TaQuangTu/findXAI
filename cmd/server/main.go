package main

import (
	"findx/internal/server"
	"findx/pkg/protogen"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		log.Fatal("GOOGLE_API_KEY environment variable not set")
	}

	engineID := os.Getenv("GOOGLE_SEARCH_ENGINE_ID")
	if engineID == "" {
		log.Fatal("GOOGLE_SEARCH_ENGINE_ID environment variable not set")
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	searchServer := server.NewSearchServer(apiKey, engineID)
	protogen.RegisterSearchServiceServer(s, searchServer)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
