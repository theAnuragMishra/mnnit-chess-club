-- +goose Up
CREATE TABLE tournaments (
    id VARCHAR(20) PRIMARY KEY CHECK (TRIM(id) <> ''),
    name VARCHAR(100) NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    duration INT NOT NULL,
    base_time INT NOT NULL,
    increment INT NOT NULL,
    status INT CHECK(status IN (0,1,2)) NOT NULL DEFAULT 0,
    berserk_allowed BOOLEAN NOT NULL DEFAULT FALSE,
    created_by INT REFERENCES users(id) ON DELETE SET NULL
);

CREATE TABLE tournament_players (
    id SERIAL PRIMARY KEY,
    player_id INT NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    tournament_id VARCHAR(20) NOT NULL REFERENCES tournaments(id) ON DELETE CASCADE,
    score INT DEFAULT 0,
    scores INT[],
    streak INT DEFAULT 0,
    UNIQUE (player_id, tournament_id)
);

CREATE TABLE games (
    id VARCHAR(20) PRIMARY KEY CHECK(TRIM(id) <> ''),
    base_time INT NOT NULL,
    increment INT NOT NULL,
    tournament_id VARCHAR(20) REFERENCES tournaments(id) ON DELETE SET NULL,
    white_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    black_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    game_length INT NOT NULL DEFAULT 0,
    result INT CHECK(result IN (0,1,2,3,4)) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    time_white INT,
    time_black INT,
    method INT NOT NULL DEFAULT 0,
    rating_w INT NOT NULL,
    rating_b INT NOT NULL,
    change_w INT,
    change_b INT,
    berserk_white BOOLEAN NOT NULL DEFAULT FALSE,
    berserk_black BOOLEAN NOT NULL DEFAULT FALSE
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
CREATE INDEX tp_pid_tid_index ON tournament_players(tournament_id, player_id);
CREATE INDEX games_w_id_index ON games(white_id);
CREATE INDEX games_b_id_index ON games(black_id);
CREATE INDEX games_tournament_id_index ON games(tournament_id);
CREATE INDEX games_created_at_index ON games(created_at DESC);
CREATE INDEX moves_game_id_index ON moves(game_id);
CREATE INDEX moves_move_number_index ON moves(move_number);

-- +goose Down
DROP INDEX IF EXISTS moves_move_number_index;
DROP INDEX IF EXISTS moves_game_id_index;
DROP INDEX IF EXISTS games_created_at_index;
DROP INDEX IF EXISTS games_tournament_id_index;
DROP INDEX IF EXISTS games_b_id_index;
DROP INDEX IF EXISTS games_w_id_index;
DROP INDEX IF EXISTS tp_pid_tid_index;
DROP TABLE IF EXISTS moves;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS tournament_players;
DROP TABLE IF EXISTS tournaments;
