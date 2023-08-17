-- psql -h localhost -U $POSTGRES_USER -f scripts/setup.sql 
-- migrate -path=./migrations/ --database=$SHORTURL_DB_DSN up

CREATE DATABASE shorturl;

\c shorturl

CREATE ROLE shorturl WITH LOGIN PASSWORD 'pa$$w0rd';
-- BEGIN: 3f7d8e5fjw9a
GRANT ALL PRIVILEGES ON DATABASE shorturl TO shorturl;
GRANT ALL PRIVILEGES ON SCHEMA public TO shorturl;

-- END: 3f7d8e5fjw9a

-- CREATE TABLE IF NOT EXISTS urls (
--     id BIGINT PRIMARY KEY,
--     long_url TEXT NOT NULL,
--     short_url TEXT NOT NULL
-- );  createdAt TIMESTAMP NOT NULL DEFAULT NOW(),
-- )