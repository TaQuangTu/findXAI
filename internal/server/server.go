package server

import (
	"context"
	"encoding/json"
	"fmt"

	"findx/internal/search"
	"findx/pkg/protogen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SearchServer struct {
	protogen.UnimplementedSearchServiceServer
	googleClient *search.Client
}

func NewSearchServer(apiKey, engineID string) *SearchServer {
	return &SearchServer{
		googleClient: search.NewClient(apiKey, engineID),
	}
}

func (s *SearchServer) Search(ctx context.Context, req *protogen.SearchRequest) (*protogen.SearchResponse, error) {
	if req.Query == "" {
		return nil, status.Error(codes.InvalidArgument, "query is required")
	}

	params := map[string]string{
		"lr":           fmt.Sprintf("lang_%s", req.Language),
		"cr":           req.Country,
		"num":          fmt.Sprintf("%d", req.NumResults),
	}

	results, err := s.googleClient.Search(ctx, req.Query, params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "search failed: %v", err)
	}

	response := &protogen.SearchResponse{
		Results: make([]*protogen.SearchResult, 0, len(results)),
	}
	fmt.Print("No da den day roi ================")
	js_res, err := json.Marshal(results)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "search failed: %v", err)
	}
	fmt.Println(string(js_res))

	for _, r := range results {
		response.Results = append(response.Results, &protogen.SearchResult{
			Title:   r.Title,
			Link:    r.Link,
			Snippet: r.Snippet,
		})
	}

	return response, nil
}
