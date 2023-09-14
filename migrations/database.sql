
-- urls table
CREATE TABLE IF NOT EXISTS urls(
    id SERIAL PRIMARY KEY,
    original VARCHAR NOT NULL,
    short_code VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_urls_original ON urls(original);
CREATE INDEX IF NOT EXISTS idx_urls_short_code ON urls(short_code);
