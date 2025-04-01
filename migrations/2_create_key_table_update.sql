-- Add NOT NULL constraints
ALTER TABLE api_keys
  ALTER COLUMN daily_queries SET NOT NULL,
  ALTER COLUMN is_active SET NOT NULL;

-- Change TIMESTAMP WITH TIME ZONE to TIMESTAMP WITHOUT TIME ZONE
ALTER TABLE api_keys
  ALTER COLUMN created_at TYPE TIMESTAMP WITHOUT TIME ZONE,
  ALTER COLUMN updated_at TYPE TIMESTAMP WITHOUT TIME ZONE;

-- Drop daily_limit and last_used columns
ALTER TABLE api_keys
  DROP COLUMN daily_limit,
  DROP COLUMN last_used;

-- Rename and update daily_queries to have default 100
ALTER TABLE api_keys
  ALTER COLUMN daily_queries SET DEFAULT 100;

-- Add new reseted_at column
ALTER TABLE api_keys
  ADD COLUMN reseted_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW();

-- Rebuild indexes (optional if you want to recreate them)
DROP INDEX IF EXISTS idx_api_keys_active;
DROP INDEX IF EXISTS idx_api_keys_usage;
CREATE INDEX idx_api_keys_active ON api_keys(is_active);
CREATE INDEX idx_api_keys_usage ON api_keys(daily_queries);
