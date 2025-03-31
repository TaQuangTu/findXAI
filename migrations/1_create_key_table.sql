CREATE TABLE api_keys (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    api_key VARCHAR(100) NOT NULL UNIQUE,
    search_engine_id VARCHAR(100) NOT NULL,
    daily_queries INTEGER NOT NULL DEFAULT 100,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    reseted_at TIMESTAMP WITHOUT TIMEZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITHOUT TIMEZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIMEZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_api_keys_active ON api_keys(is_active);
CREATE INDEX idx_api_keys_usage ON api_keys(daily_queries);


-- INSERT INTO api_keys (
--     name, 
--     api_key, 
--     search_engine_id,
--     daily_limit,
--     is_active
-- ) VALUES (
--     'key-name',      -- Unique name for the key
--     'your api key',  -- Google API key
--     'search engine id',    -- Search engine ID
--     100,                     -- Daily query limit
--     TRUE                     -- Active status
-- )
-- RETURNING id, created_at;