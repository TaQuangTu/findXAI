ALTER TABLE api_keys
  ADD COLUMN IF NOT EXISTS status_code INTEGER NOT NULL DEFAULT 0;

ALTER TABLE api_keys
  ADD COLUMN IF NOT EXISTS error_msg VARCHAR(1024) NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_api_keys_status_code ON api_keys(status_code);

CREATE INDEX IF NOT EXISTS idx_api_keys_active_queries ON api_keys(is_active, daily_queries);

CREATE INDEX IF NOT EXISTS idx_api_keys_api_key ON api_keys(api_key);
