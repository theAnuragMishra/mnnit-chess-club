-- +goose Up
CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    username VARCHAR(20) UNIQUE,
    avatar_url TEXT,
    google_id VARCHAR(255) UNIQUE NOT NULL,
    rating FLOAT NOT NULL DEFAULT 1500.0,
    rd FLOAT NOT NULL DEFAULT 350.0,
    volatility FLOAT NOT NULL DEFAULT 0.06
);

-- +goose Down
DROP TABLE users;
