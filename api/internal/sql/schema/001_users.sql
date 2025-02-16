-- +goose Up
CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    username VARCHAR(20) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);




-- +goose Down
DROP TABLE users;
