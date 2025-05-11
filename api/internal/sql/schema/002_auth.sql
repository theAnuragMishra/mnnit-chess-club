-- +goose Up
CREATE TABLE sessions(
                         id TEXT PRIMARY KEY,
                         user_id INTEGER NOT NULL,
                         created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         expires_at TIMESTAMPTZ NOT NULL,
                         FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX sessions_user_id_index ON sessions(user_id);

-- +goose Down
DROP INDEX IF EXISTS sessions_user_id_index;
DROP TABLE sessions;
