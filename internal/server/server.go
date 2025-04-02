package server

import (
	"context"
	"findx/config"
	"findx/internal/lockdb"
	"findx/internal/search"
	"findx/pkg/protogen"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SearchServer struct {
	protogen.UnimplementedSearchServiceServer
	KeyManager   *search.ApiKeyManager
	googleClient *search.Client

	lockDb      lockdb.ILockDb
	rateLimiter lockdb.RateLimiter
}

func NewSearchServer(conf *config.Config, lockDb lockdb.ILockDb, rateLimiter lockdb.RateLimiter) *SearchServer {
	return &SearchServer{
		googleClient: search.NewClient(),
		KeyManager:   search.NewApiKeyManager(conf.POSTGRES_DSN, lockDb, rateLimiter),
		lockDb:       lockDb,
		rateLimiter:  rateLimiter,
	}
}

func (s *SearchServer) Search(ctx context.Context, req *protogen.SearchRequest) (*protogen.SearchResponse, error) {
	if req.Query == "" {
		return nil, status.Error(codes.InvalidArgument, "query is required")
	}

	// Number of bucket can be configured
	bucketList, err := s.KeyManager.GetKeyBucket(ctx, 5)
	if err != nil {
		return nil, status.Error(codes.ResourceExhausted, err.Error())
	}
	availableKey, err := s.KeyManager.GetAvailableKey(ctx, bucketList)
	if err != nil {
		return nil, status.Error(codes.ResourceExhausted, err.Error())
	}

	var (
		weShouldDoSomething bool
	)
	if weShouldDoSomething = bucketList.Avg() < 10; weShouldDoSomething {
		fmt.Println("nooooooooo")
	}

	params := map[string]string{
		"lr":  fmt.Sprintf("lang_%s", req.Language),
		"cr":  req.Country,
		"num": fmt.Sprintf("%d", req.NumResults),
	}

	results, err := s.googleClient.Search(ctx, availableKey.ApiKey, availableKey.EngineId, req.Query, params)
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
