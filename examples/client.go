package main

import (
	"context"
	"log"
	"time"

	"findx/pkg/contentsvc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := contentsvc.NewContentServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := c.ExtractContentFromLinks(ctx, &contentsvc.ExtractContentFromLinksRequest{
		Links: []string{"https://vov.vn/thi-truong/chung-khoan/chung-khoan-sang-44-vn-index-tiep-tuc-lao-doc-boc-hoi-gan-50-diem-post1189618.vov"},
	})
	if err != nil {
		log.Fatalf("could not search: %v", err)
	}

	for _, result := range resp.Contents {
		log.Printf("Title: %s", result.Title)
		log.Printf("Content: %s", result.Content)
		log.Printf("Link: %s", result.Link)
		log.Println("------------------------------")
	}
}
