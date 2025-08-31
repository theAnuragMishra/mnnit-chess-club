-- +goose Up
CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    role INT NOT NULL DEFAULT 0,
    username VARCHAR(20) UNIQUE,
    avatar_url TEXT,
    google_id VARCHAR(255) UNIQUE NOT NULL,
    rating FLOAT NOT NULL DEFAULT 1500.0,
    rd FLOAT NOT NULL DEFAULT 350.0,
    volatility FLOAT NOT NULL DEFAULT 0.06
);

CREATE INDEX users_username_index ON users(username);
CREATE INDEX users_rating_index ON users(rating DESC);

-- +goose Down
DROP INDEX IF EXISTS users_rating_index;
DROP INDEX IF EXISTS users_username_index;
DROP TABLE users;
