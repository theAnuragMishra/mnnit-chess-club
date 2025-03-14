-- +goose Up
CREATE TABLE games (
    id SERIAL PRIMARY KEY,
    base_time INT NOT NULL,
    increment INT NOT NULL,
    white_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    black_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    white_username VARCHAR(20) REFERENCES users(username) ON DELETE SET NULL,
    black_username VARCHAR(20) REFERENCES users(username) ON DELETE SET NULL,
    fen TEXT NOT NULL,
    game_length SMALLINT NOT NULL DEFAULT 0,
    result VARCHAR(10) CHECK(result IN ('1-0', '0-1', '1/2-1/2', 'ongoing')) NOT NULL DEFAULT 'ongoing',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    end_time_left_white INT,
    end_time_left_black INT,
    result_reason VARCHAR(100)
);

CREATE TABLE moves (
    id SERIAL PRIMARY KEY ,
    game_id INTEGER REFERENCES games(id) ON DELETE CASCADE NOT NULL,
    move_number INT NOT NULL,
    player_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    move_notation VARCHAR(10) NOT NULL,
    orig VARCHAR(4) NOT NULL,
    dest VARCHAR(4) NOT NULL,
    move_fen TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE moves;
DROP TABLE games;
