-- +goose Up
CREATE TABLE tournaments (
    id VARCHAR(20) PRIMARY KEY CHECK (TRIM(id) <> ''),
    name VARCHAR(100) NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    duration INT NOT NULL,
    base_time INT NOT NULL,
    increment INT NOT NULL,
    created_by INT REFERENCES users(id) ON DELETE SET NULL
);

CREATE TABLE tournament_players (
    id SERIAL PRIMARY KEY,
    player_id INT NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    tournament_id VARCHAR(20) NOT NULL REFERENCES tournaments(id) ON DELETE CASCADE,
    score INT DEFAULT 0,
    scores SMALLINT[],
    UNIQUE (player_id, tournament_id)
);

CREATE TABLE games (
    id VARCHAR(20) PRIMARY KEY CHECK(TRIM(id) <> ''),
    base_time INT NOT NULL,
    increment INT NOT NULL,
    tournament_id VARCHAR(20) REFERENCES tournaments(id) ON DELETE SET NULL,
    white_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    black_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    game_length SMALLINT NOT NULL DEFAULT 0,
    result VARCHAR(10) CHECK(result IN ('1-0', '0-1', '1/2-1/2', 'ongoing', 'aborted')) NOT NULL DEFAULT 'ongoing',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end_time_left_white INT,
    end_time_left_black INT,
    result_reason VARCHAR(100),
    rating_w INT NOT NULL,
    rating_b INT NOT NULL,
    change_w INT,
    change_b INT
);

CREATE TABLE moves (
    id SERIAL PRIMARY KEY ,
    game_id VARCHAR(20) REFERENCES games(id) ON DELETE CASCADE NOT NULL,
    move_number INT NOT NULL,
    move_notation VARCHAR(10) NOT NULL,
    orig VARCHAR(4) NOT NULL,
    dest VARCHAR(4) NOT NULL,
    move_fen TEXT NOT NULL,
    time_left INT
);
CREATE INDEX games_w_id_index ON games(white_id);
CREATE INDEX games_b_id_index ON games(black_id);
CREATE INDEX games_wb_id_index ON games(white_id, black_id);
CREATE INDEX games_result_index ON games(result);
CREATE INDEX games_created_at_index ON games(created_at);
CREATE INDEX moves_game_id_index ON moves(game_id);
CREATE INDEX moves_move_number_index ON moves(move_number);

-- +goose Down
DROP INDEX IF EXISTS moves_move_number_index;
DROP INDEX IF EXISTS moves_game_id_index;
DROP INDEX IF EXISTS games_created_at_index;
DROP INDEX IF EXISTS games_result_index;
DROP INDEX IF EXISTS games_wb_id_index;
DROP INDEX IF EXISTS games_b_id_index;
DROP INDEX IF EXISTS games_w_id_index;
DROP TABLE IF EXISTS moves;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS tournament_players;
DROP TABLE IF EXISTS tournaments;
