package search

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

type ApiKeyManager struct {
	db *sql.DB
}

func NewApiKeyManager(dsn string) *ApiKeyManager {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("failed to connect to database: %w", err)
		panic(err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	if err := db.Ping(); err != nil {
		panic(fmt.Errorf("failed to ping database: %w", err))
	}

	return &ApiKeyManager{db: db}
}

func (m *ApiKeyManager) GetAvailableKey(ctx context.Context) (apiKey, engineID string, err error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return "", "", err
	}
	defer tx.Rollback()

	// Find the first available key that hasn't reached its daily limit
	row := tx.QueryRowContext(ctx, `
		UPDATE api_keys 
		SET daily_queries = daily_queries + 1, 
		    last_used = NOW() 
		WHERE id = (
			SELECT id 
			FROM api_keys 
			WHERE is_active = TRUE AND daily_queries < daily_limit 
			ORDER BY last_used NULLS FIRST 
			FOR UPDATE SKIP LOCKED 
			LIMIT 1
		) 
		RETURNING api_key, search_engine_id
	`)

	err = row.Scan(&apiKey, &engineID)
	if err != nil {
		return "", "", err
	}

	return apiKey, engineID, tx.Commit()
}

func (m *ApiKeyManager) ResetDailyCounts() {
	// Run this at midnight UTC
	_, _ = m.db.Exec(`UPDATE api_keys SET daily_queries = 0`)
}
