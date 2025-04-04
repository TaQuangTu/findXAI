package server

import (
	"context"
	"database/sql"
	"findx/config"
	"findx/internal/liberror"
	"findx/internal/lockdb"
	"findx/internal/search"
	"findx/internal/system"
	"findx/pkg/searchsvc"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SearchServer struct {
	searchsvc.UnimplementedSearchServiceServer
	KeyManager   *search.ApiKeyManager
	googleClient *search.Client

	db          *sql.DB
	lockDb      lockdb.ILockDb
	rateLimiter lockdb.RateLimiter
}

func NewSearchServer(conf *config.Config, lockDb lockdb.ILockDb, rateLimiter lockdb.RateLimiter) *SearchServer {
	db, err := sql.Open("postgres", conf.POSTGRES_DSN)
	system.PanicOnError(err)
	system.RegisterRootCloser(db.Close)

	// Configure connection pool
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	err = db.Ping()
	system.PanicOnError(err)
	return &SearchServer{
		db:           db,
		googleClient: search.NewClient(),
		KeyManager:   search.NewApiKeyManager(conf, db, lockDb, rateLimiter),
		lockDb:       lockDb,
		rateLimiter:  rateLimiter,
	}
}

func (s *SearchServer) Search(ctx context.Context, req *searchsvc.SearchRequest) (*searchsvc.SearchResponse, error) {
	if err := ValidateSearchRequest(req); err != nil {
		return nil, Error(
			liberror.WrapStack(err, "search: user use invalid parameter"),
		)
	}

	if req.Q == "" {
		return nil, status.Error(codes.InvalidArgument, "query is required")
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, Error(
			liberror.WrapStack(err, "search: tx begin: failed"),
		)
	}
	defer tx.Rollback()

	availableKey, err := s.KeyManager.GetAvailableKey(ctx)
	if err != nil {
		return nil, Error(err)
	}

	// convert req to map[string]string
	params := ProtoMessageToMap(req)

	results, statusCode, searchErr := s.googleClient.Search(ctx, availableKey.ApiKey, availableKey.EngineId, params)

	var msg string
	if searchErr != nil {
		msg = searchErr.Error()
	}
	err = s.KeyManager.UpdateKeyStatus(ctx, tx, availableKey.ApiKey, statusCode, msg)
	if err != nil {
		return nil, Error(err)
	}
	err = tx.Commit()
	if err != nil {
		return nil, Error(
			liberror.WrapStack(err, "search: tx commit: failed"),
		)
	}

	if searchErr != nil {
		return nil, Error(
			liberror.WrapStack(err, "search: failed"),
		)
	}
	response := &searchsvc.SearchResponse{
		Results: make([]*searchsvc.SearchResult, 0, len(results)),
	}
	for _, r := range results {
		response.Results = append(response.Results, &searchsvc.SearchResult{
			Title:   r.Title,
			Link:    r.Link,
			Snippet: r.Snippet,
		})
	}
	return response, nil
}

func (s *SearchServer) DeactivateKeys(ctx context.Context, req *searchsvc.DeactivateKeysRequest) (resModel *searchsvc.DeactivateKeysResponse, err error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		err = liberror.WrapStack(err, "deactivate key: begin tx failed")
		return
	}
	if len(req.ApiKeys) <= 0 {
		err = liberror.WrapStack(liberror.ErrorDataInvalid, "list key is empty")
		return
	}
	if req.ForceDelete {
		err = Error(s.KeyManager.HardDeleteKeys(ctx, tx, req.ApiKeys))
		return
	}
	err = Error(s.KeyManager.UpdateKeyActiveStatus(ctx, tx, req.ApiKeys, false))
	return
}

func (s *SearchServer) ActivateKeys(ctx context.Context, req *searchsvc.ActivateKeysRequest) (resModel *searchsvc.ActivateKeysResponse, err error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		err = liberror.WrapStack(err, "activate key: begin tx failed")
		return
	}
	if len(req.ApiKeys) <= 0 {
		err = liberror.WrapStack(liberror.ErrorDataInvalid, "list key is empty")
		return
	}
	err = Error(s.KeyManager.UpdateKeyActiveStatus(ctx, tx, req.ApiKeys, true))
	return
}

func (s *SearchServer) AddKeys(ctx context.Context, req *searchsvc.AddKeysRequest) (resModel *searchsvc.AddKeysResponse, err error) {
	var (
		sqlInput = make([][]any, 0)
	)
	for idx, item := range req.Data {
		if item.ApiKey == "" || item.Name == "" || item.SearchEngineId == "" {
			err = Error(liberror.
				WrapStack(liberror.ErrorDataInvalid, "add key: invalid data").
				WithField("item_idx", idx))
			return
		}
		sqlInput = append(sqlInput, []any{item.Name, item.ApiKey, item.SearchEngineId})
	}
	err = Error(s.KeyManager.AddKeys(ctx, sqlInput))
	return
}

func (s *SearchServer) GetKeys(ctx context.Context, req *searchsvc.GetKeysRequest) (_ *searchsvc.GetKeysResponse, err error) {
	if len(req.ApiKeys) == 0 {
		err = Error(
			liberror.WrapStack(liberror.ErrorDataInvalid, "list key is empty"))
		return
	}
	keys, err := s.KeyManager.GetKeys(ctx, req.ApiKeys)
	if err != nil {
		err = Error(err)
		return
	}
	var (
		resModel = &searchsvc.GetKeysResponse{
			Results: make([]*searchsvc.KeyInfo, len(keys)),
		}
	)
	for idx, key := range keys {
		resModel.Results[idx] = &searchsvc.KeyInfo{
			Id:             key.Id,
			Name:           key.Name,
			ApiKey:         key.ApiKey,
			SearchEngineId: key.SearchEngineId,
			IsActive:       key.IsActive,
			DailyQueries:   key.DailyQueries,
			StatusCode:     key.StatusCode,
			ErrorMsg:       key.ErrorMsg,
			CreatedAt:      key.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      key.UpdatedAt.Format(time.RFC3339),
		}
	}
	return resModel, nil
}
