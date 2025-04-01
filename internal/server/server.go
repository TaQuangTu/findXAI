package server

import (
	"context"
	"findx/config"
	"findx/internal/lockdb"
	"findx/internal/search"
	"findx/pkg/protogen"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SearchServer struct {
	protogen.UnimplementedSearchServiceServer
	keyManager   *search.ApiKeyManager
	googleClient *search.Client

	lockDb      lockdb.ILockDb
	rateLimiter lockdb.RateLimiter
}

func NewSearchServer(conf *config.Config, lockDb lockdb.ILockDb, rateLimiter lockdb.RateLimiter) *SearchServer {
	return &SearchServer{
		googleClient: search.NewClient(),
		keyManager:   search.NewApiKeyManager(conf.POSTGRES_DSN, lockDb, rateLimiter),
		lockDb:       lockDb,
		rateLimiter:  rateLimiter,
	}
}

func (s *SearchServer) Search(ctx context.Context, req *protogen.SearchRequest) (*protogen.SearchResponse, error) {
	if req.Query == "" {
		return nil, status.Error(codes.InvalidArgument, "query is required")
	}

	// Number of bucket can be configured
	bucketList, err := s.keyManager.GetKeyBucket(ctx, 5)
	if err != nil {
		return nil, status.Errorf(codes.ResourceExhausted, "error")
	}
	availableKey, err := s.keyManager.GetAvailableKey(ctx, bucketList)
	if err != nil {
		return nil, status.Errorf(codes.ResourceExhausted, "no available API keys")
	}
	defer func() {
		var (
			dateOnlyCurrentTime = time.Now().UTC().Truncate(24 * time.Hour)
			dateOnlyUpdatedTime = availableKey.ResetedAt.UTC().Truncate(24 * time.Hour)
		)
		if !dateOnlyCurrentTime.
			After(dateOnlyUpdatedTime) {
			return
		}
		goCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		go func(ctx context.Context) {
			//TODO: handle error
			ourLock, err := s.lockDb.LockSimple(ctx, "search:get_key:reset")
			if err != nil {
				return
			}
			defer ourLock.Unlock()
			s.keyManager.ResetDailyCounts(100, availableKey.ResetedAt)
		}(goCtx)
	}()

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
