package search

import (
	"context"
	"database/sql"
	"findx/config"
	"findx/internal/liberror"
	"findx/internal/lockdb"
	"fmt"
	"github.com/lib/pq"
	"log"
	"time"

	"github.com/go-redis/redis_rate/v9"
)

type ApiKeyManager struct {
	db          *sql.DB
	lockDb      lockdb.ILockDb
	rateLimiter lockdb.RateLimiter

	conf *config.Config
}

func NewApiKeyManager(conf *config.Config, db *sql.DB, lockDb lockdb.ILockDb, rateLimiter lockdb.RateLimiter) *ApiKeyManager {
	return &ApiKeyManager{
		conf:        conf,
		db:          db,
		lockDb:      lockDb,
		rateLimiter: rateLimiter}
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
	bucketList, err := m.GetKeyBucket(ctx, m.conf.APP_KEY_BUCKET)
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
	res, err := m.db.Exec(`
		UPDATE api_keys
			SET daily_queries = $1 , reseted_at = NOW()
		WHERE reseted_at::DATE < CURRENT_DATE
	`, limit)
	if err != nil {
		err = liberror.WrapStack(err, "reset daily counts: update failed")
	}
	affectedRows, err := res.RowsAffected()
	if err != nil {
		err = liberror.WrapStack(err, "reset daily counts: get rows affected failed")
	}
	log.Println("[SUCCESS]: Reset daily counts for ", affectedRows, " keys")
}

func (m *ApiKeyManager) UpdateKeyStatus(ctx context.Context, apiKey string, status int, msg string) (_ error) {
	_, err := m.db.Exec(`
		UPDATE api_keys
			SET status_code = $1, error_msg = $2, updated_at = NOW()
		WHERE api_key = $3
	`, &status, &msg, &apiKey)
	if err != nil {
		err = liberror.WrapStack(err, "key status: update failed")
	}
	return err
}

func (m *ApiKeyManager) UpdateKeyActiveStatus(ctx context.Context, apiKeys []string, isActivate bool) (_ error) {
	_, err := m.db.Exec(`
		UPDATE api_keys
			SET is_active = $1
		WHERE api_key = ANY($2::text[])
	`, isActivate, pq.Array(apiKeys))
	if err != nil {
		err = liberror.WrapStack(err, "key active status: update faield").
			WithFields(liberror.AdditionalData{
				"api_keys":   apiKeys,
				"isActivate": isActivate,
			})
	}
	return err
}

func (m *ApiKeyManager) HardDeleteKeys(ctx context.Context, apiKeys []string) (_ error) {
	_, err := m.db.Exec(`
		DELETE FROM api_keys
		WHERE api_key = ANY($1::text[])
	`, pq.Array(apiKeys))
	if err != nil {
		err = liberror.WrapStack(err, "api keys: hard delete faield").
			WithFields(liberror.AdditionalData{
				"api_keys": apiKeys,
			})
	}
	return err
}

func (m *ApiKeyManager) GetKeys(ctx context.Context, apiKeys []string) (keys []Key, err error) {
	rows, err := m.db.QueryContext(ctx, `
			SELECT
				id,
				name,
				api_key,
				search_engine_id,
				is_active,
				daily_queries,
				status_code,
				error_msg,
				created_at,
				updated_at
			FROM api_keys
			WHERE api_key = ANY($1::text[])
		`, pq.Array(apiKeys))
	defer rows.Close()

	for rows.Next() {
		var (
			key Key
		)
		err = rows.Scan(&key.Id, &key.Name, &key.ApiKey, &key.IsActive, &key.DailyQueries, &key.StatusCode, &key.ErrorMsg, &key.CreatedAt, &key.UpdatedAt)
		if err != nil {
			err = liberror.WrapStack(err, "get key: scan data failed")
			return
		}
		keys = append(keys, key)
	}
	return keys, nil
}

func (m *ApiKeyManager) AddKeys(ctx context.Context, data [][]any) (err error) {
	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		err = liberror.WrapStack(err, "add key: begin tx failed")
		return
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO api_keys (name, api_key, search_engine_id)
			VALUES ($1, $2, $3)
	`)
	if err != nil {
		err = liberror.WrapStack(err, "add key: preapre stmt failed")
		return
	}
	defer stmt.Close()

	for _, row := range data {
		_, err = stmt.Exec(row...)
		if err != nil {
			_ = tx.Rollback()
			err = liberror.WrapStack(err, "add key: exec failed").
				WithField("data", row)
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		err = liberror.WrapStack(err, "add key: commit tx failed")
	}
	return
}
