-- +goose Up
CREATE TABLE games (
    id SERIAL PRIMARY KEY,
    white_player_id UUID REFERENCES users(id) ON DELETE SET NULL,
    black_player_id UUID REFERENCES users(id) ON DELETE SET NULL,
    result VARCHAR(10) CHECK(result IN ('1-0', '0-1', '1/2-1/2', 'ongoing')) DEFAULT 'ongoing',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP NULL
);

CREATE TABLE moves (
    id SERIAL PRIMARY KEY ,
    game_id INT REFERENCES games(id) ON DELETE CASCADE NOT NULL,
    move_number INT NOT NULL,
    player_id UUID REFERENCES users(id) ON DELETE CASCADE,
    move_notation VARCHAR(10) NOT NULL,
    move_fen TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE games;
DROP TABLE moves;