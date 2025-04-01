package search

import (
	"context"
	"database/sql"
	"findx/internal/lockdb"
	"findx/internal/system"
	"fmt"
	_ "github.com/lib/pq"
	"time"

	"github.com/go-redis/redis_rate/v9"
	_ "github.com/lib/pq"
)

type ApiKeyManager struct {
	db          *sql.DB
	lockDb      lockdb.ILockDb
	rateLimiter lockdb.RateLimiter
}

func NewApiKeyManager(dsn string, lockDb lockdb.ILockDb, rateLimiter lockdb.RateLimiter) *ApiKeyManager {
	db, err := sql.Open("postgres", dsn)
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
	if err := db.Ping(); err != nil {
		panic(fmt.Errorf("failed to ping database: %w", err))
	}

	return &ApiKeyManager{db: db, lockDb: lockDb, rateLimiter: rateLimiter}
}

func (m *ApiKeyManager) GetKeyBucket(ctx context.Context, numberOfPartition int) (_ KeyBucketList, err error) {
	rows, err := m.db.QueryContext(ctx, `
		SELECT
				AVG(daily_queries) AS partition_avg,
				(id % ?) AS partition_id
		FROM api_keys
		GROUP BY partition_id
		ORDER BY partition_avg DESC
`, numberOfPartition)
	if err != nil {
		return
	}
	defer rows.Close()
	var (
		bucketList = make([]KeyBucket, 0)
	)

	for rows.Next() {
		var (
			partitionId     int
			partitionAvg    float64
			numberOfRecords int
		)
		err = rows.Scan(&partitionAvg, &numberOfRecords, &partitionId)
		if err != nil {
			return
		}
		bucketList = append(bucketList, KeyBucket{
			NumberOfPartition: numberOfPartition,
			PartitionId:       partitionId,
			PartitionAvg:      partitionAvg,
		})
	}
	if len(bucketList) <= 0 {
		err = fmt.Errorf("invalid param")
		return
	}
	// First bucket will always has highest capacity
	// Some last buckets are possible running out of limit
	return bucketList, err
}

func (m *ApiKeyManager) GetAvailableKey(ctx context.Context, bucketList KeyBucketList) (_ *AvailableKey, err error) {
	var (
		selectedBucket KeyBucket
	)
	for _, bucket := range bucketList {
		var (
			key       = fmt.Sprintf("search:get_key:partition:%d", bucket.PartitionId)
			rateLimit = int(bucket.PartitionAvg)
		)
		if rateLimit < 1 {
			rateLimit = 1
		}
		result, err := m.rateLimiter.Allow(ctx, key, redis_rate.Limit{
			Rate:   rateLimit,
			Burst:  rateLimit,
			Period: time.Second,
		})
		if err != nil {
			return nil, err
		}
		if result.Allowed > 0 {
			selectedBucket = bucket
			break
		} else {
			continue
		}
	}

	if selectedBucket.PartitionAvg <= 0 {
		return nil, fmt.Errorf("out of limit")
	}
	// Find the first available key that hasn't reached its daily limit in highest capacity bucket
	row := m.db.QueryRowContext(ctx, `
		WITH selected_key AS (
    		SELECT id
    		FROM api_keys
    		WHERE
        		is_active = TRUE
        		AND id % ? = ?
        		AND daily_queries > 0
    		ORDER BY daily_queries DESC
    		FOR UPDATE SKIP LOCKED
    		LIMIT 1
		)
		UPDATE api_keys
		SET
    		daily_queries = daily_queries - 1,
    		updated_at = NOW()
		WHERE id IN (SELECT id FROM selected_key)
		RETURNING api_key, search_engine_id, reseted_at;
	`, selectedBucket.NumberOfPartition, selectedBucket.PartitionId)

	var availableKey AvailableKey
	err = row.Scan(&availableKey.ApiKey, &availableKey.EngineId, availableKey.ResetedAt)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("out of limit")
	}

	return &availableKey, nil
}

func (m *ApiKeyManager) ResetDailyCounts(limit int, updatedAt time.Time) {
	// Run this at midnight UTC
	_, _ = m.db.Exec(`
		UPDATE api_keys
			SET daily_queries = ? , reseted_at = NOW()
		WHERE
			reseted_at::DATE < ?
	`, limit, updatedAt.Format("2006-01-02"))
}
