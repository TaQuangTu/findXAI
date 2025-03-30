package main

import (
	"context"
	"log"
	"time"

	"findx/pkg/protogen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := protogen.NewSearchServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := c.Search(ctx, &protogen.SearchRequest{
		Query:      "Thủ Tướng Phạm Minh Chính đang ở đâu",
		Language:   "vi",
		Country:    "vi",
		NumResults: 10,
		// StartDate:  "2023-01-01",
		// EndDate:    "2023-12-31",
	})
	if err != nil {
		log.Fatalf("could not search: %v", err)
	}

	log.Printf("Results:")
	for _, result := range resp.Results {
		log.Printf("- %s (%s)", result.Title, result.Link)
		log.Printf("  %s", result.Snippet)
	}
}
