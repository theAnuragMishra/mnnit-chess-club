-- +goose Up

CREATE TABLE sessions(
                         id TEXT PRIMARY KEY,
                         user_id INTEGER NOT NULL,
                         created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         expires_at TIMESTAMPTZ NOT NULL,
                         FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE csrf_tokens(
                            session_id TEXT PRIMARY KEY ,
                            token TEXT NOT NULL,
                            created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            expires_at TIMESTAMPTZ NOT NULL,
                            FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

-- +goose Down

DROP TABLE csrf_tokens;
DROP TABLE sessions;
