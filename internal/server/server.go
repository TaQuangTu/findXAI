package server

import (
	"context"
	"fmt"

	"findx/internal/search"
	"findx/pkg/protogen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SearchServer struct {
	protogen.UnimplementedSearchServiceServer
	keyManager   *search.ApiKeyManager
	googleClient *search.Client
}

func NewSearchServer(dsn string) *SearchServer {
	return &SearchServer{
		googleClient: search.NewClient(),
		keyManager:   search.NewApiKeyManager(dsn),
	}
}

func (s *SearchServer) Search(ctx context.Context, req *protogen.SearchRequest) (*protogen.SearchResponse, error) {
	if req.Query == "" {
		return nil, status.Error(codes.InvalidArgument, "query is required")
	}

	apiKey, engineID, err := s.keyManager.GetAvailableKey(ctx)
	if err != nil {
		return nil, status.Errorf(codes.ResourceExhausted, "no available API keys")
	}

	params := map[string]string{
		"lr":  fmt.Sprintf("lang_%s", req.Language),
		"cr":  req.Country,
		"num": fmt.Sprintf("%d", req.NumResults),
	}

	results, err := s.googleClient.Search(ctx, apiKey, engineID, req.Query, params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "search failed: %v", err)
	}

	response := &protogen.SearchResponse{
		Results: make([]*protogen.SearchResult, 0, len(results)),
	}

	for _, r := range results {
		response.Results = append(response.Results, &protogen.SearchResult{
			Title:   r.Title,
			Link:    r.Link,
			Snippet: r.Snippet,
		})
	}

	return response, nil
}
