CREATE TABLE IF NOT EXISTS urls (
    id BIGINT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    long_url TEXT NOT NULL,
    short_url TEXT NOT NULL,
    visits INT NOT NULL DEFAULT 0
);