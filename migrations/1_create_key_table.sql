CREATE TABLE api_keys (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    api_key VARCHAR(100) NOT NULL UNIQUE,
    search_engine_id VARCHAR(100) NOT NULL,
    daily_queries INTEGER DEFAULT 0,
    daily_limit INTEGER DEFAULT 100,
    last_used TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
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
