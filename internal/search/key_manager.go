package search

import (
	"context"
	"database/sql"
	"findx/internal/liberror"
	"findx/internal/lockdb"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/go-redis/redis_rate/v9"
	_ "github.com/lib/pq"
)

type ApiKeyManager struct {
	db          *sql.DB
	lockDb      lockdb.ILockDb
	rateLimiter lockdb.RateLimiter
}

func NewApiKeyManager(db *sql.DB, lockDb lockdb.ILockDb, rateLimiter lockdb.RateLimiter) *ApiKeyManager {

	return &ApiKeyManager{db: db, lockDb: lockDb, rateLimiter: rateLimiter}
}

func (m *ApiKeyManager) GetKeyBucket(ctx context.Context, numberOfPartition int) (_ KeyBucketList, err error) {
	rows, err := m.db.QueryContext(ctx, `
		SELECT
				AVG(daily_queries) AS partition_avg,
				COUNT(*) AS record_count,
				(id % $1) AS partition_id
		FROM api_keys
		GROUP BY partition_id
		ORDER BY partition_avg DESC
	`, numberOfPartition)
	if err != nil {
		err = liberror.WrapStack(err, "get bucket: query failed")
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
			err = liberror.WrapStack(err, "get bucket: scan failed")
			return
		}
		bucketList = append(bucketList, KeyBucket{
			NumberOfPartition: numberOfPartition,
			PartitionId:       partitionId,
			PartitionAvg:      partitionAvg,
		})
	}
	if len(bucketList) <= 0 {
		err = liberror.WrapStack(liberror.ErrorNotFound, "get bucket: no bucket found")
		return
	}
	// First bucket will always has highest capacity
	// Some last buckets are possible running out of limit
	return bucketList, err
}

func (m *ApiKeyManager) GetAvailableKey(ctx context.Context) (_ *AvailableKey, err error) {
	bucketList, err := m.GetKeyBucket(ctx, 3)
	if err != nil {
		return
	}
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
			err = liberror.WrapStack(err, "get available key: rate lock failed")
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
		err = liberror.WrapStack(liberror.ErrorLimitReached, "get available key: out of capacity")
		return
	}
	// Find the first available key that hasn't reached its daily limit in highest capacity bucket
	row := m.db.QueryRowContext(ctx, `
		WITH selected_key AS (
    		SELECT id
    		FROM api_keys
    		WHERE
        		is_active = TRUE
        		AND id % $1 = $2
        		AND daily_queries > 0
						AND status_code BETWEEN 200 AND 299
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
	err = row.Scan(&availableKey.ApiKey, &availableKey.EngineId, &availableKey.ResetedAt)
	if err == nil {
		return &availableKey, nil
	}
	if err == sql.ErrNoRows {
		err = liberror.ErrorLimitReached
	}
	err = liberror.WrapStack(err, "get available key: get key failed")
	return
}

func (m *ApiKeyManager) ResetDailyCounts(limit int) {
	// Run this at midnight UTC
	_, _ = m.db.Exec(`
		UPDATE api_keys
			SET daily_queries = $1 , reseted_at = NOW()
		WHERE
			reseted_at::DATE < NOW()::DATE
	`, limit)
}

func (m *ApiKeyManager) UpdateKeyStatus(ctx context.Context, apiKey string, status int, msg string) (_ error) {
	_, err := m.db.Exec(`
		UPDATE api_keys
			SET status_code = $1, error_msg = $2
		WHERE api_key = $3
	`, &status, &msg, &apiKey)
	if err != nil {
		err = liberror.WrapStack(err, "key status: update failed")
	}
	return err
}
