package server

import (
	"context"
	"database/sql"
	"findx/config"
	"findx/internal/liberror"
	"findx/internal/lockdb"
	"findx/internal/search"
	"findx/internal/system"
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

	db          *sql.DB
	lockDb      lockdb.ILockDb
	rateLimiter lockdb.RateLimiter
}

func NewSearchServer(conf *config.Config, lockDb lockdb.ILockDb, rateLimiter lockdb.RateLimiter) *SearchServer {
	db, err := sql.Open("postgres", conf.POSTGRES_DSN)
	if err != nil {
		fmt.Println("failed to connect to database: %w", err)
		panic(err)
	}
	system.RegisterRootCloser(db.Close)

	// Configure connection pool
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	system.PanicOnError(fmt.Errorf("ping db faild: %w", db.Ping()))
	return &SearchServer{
		googleClient: search.NewClient(),
		keyManager:   search.NewApiKeyManager(db, lockDb, rateLimiter),
		db:           db,
		lockDb:       lockDb,
		rateLimiter:  rateLimiter,
	}
}

// TODO: format later
func (s *SearchServer) toError(err error) error {
	return err
}

func (s *SearchServer) Search(ctx context.Context, req *protogen.SearchRequest) (*protogen.SearchResponse, error) {
	if req.Query == "" {
		return nil, status.Error(codes.InvalidArgument, "query is required")
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, s.toError(
			liberror.WrapStack(err, "search: tx begin: failed"),
		)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	availableKey, err := s.keyManager.GetAvailableKey(ctx)
	if err != nil {
		return nil, s.toError(err)
	}
	params := map[string]string{
		"lr":  fmt.Sprintf("lang_%s", req.Language),
		"cr":  req.Country,
		"num": fmt.Sprintf("%d", req.NumResults),
	}
	results, statusCode, searchErr := s.googleClient.Search(ctx, availableKey.ApiKey, availableKey.EngineId, req.Query, params)

	var msg string
	if searchErr != nil {
		msg = searchErr.Error()
	}
	err = s.keyManager.UpdateKeyStatus(ctx, availableKey.ApiKey, statusCode, msg)
	if err != nil {
		return nil, s.toError(err)
	}
	err = tx.Commit()
	if err != nil {
		return nil, s.toError(
			liberror.WrapStack(err, "search: tx commit: failed"),
		)
	}

	if searchErr != nil {
		return nil, s.toError(
			liberror.WrapStack(err, "search: failed"),
		)
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
