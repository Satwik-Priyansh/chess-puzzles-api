-- Migration script to initialize the core database schema
BEGIN;

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

--Table 1: Users
CREATE TABLE IF NOT EXISTS users (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 email text UNIQUE NOT NULL,
 password_hash text, --bycrypt
 username text UNIQUE NOT NULL,
 rating float DEFAULT 1000,
 rating_deviation float DEFAULT 350,--high RD="unproven" new user
 created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP --timezone aware
);
-- Table 2: Puzzles
CREATE TABLE IF NOT EXISTS puzzles(
    id text PRIMARY KEY,
    fen text,
    moves text[], --split using spaces
    rating float,
    rating_deviation float,
    popularity int,
    nb_plays int,
    themes text[],  --space seprated
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--Table 3:Solve_attempts
CREATE TABLE IF NOT EXISTS solve_attempts(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    puzzle_id text NOT NULL,
    success BOOLEAN NOT NULL,
    puzzle_rating_before FLOAT NOT NULL,
    puzzle_rating_after FLOAT NOT NULL,
    user_rating_before FLOAT NOT NULL,
    user_rating_after FLOAT NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,

    --Foreign Key Declaration
    CONSTRAINT fk_attempt_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,

    CONSTRAINT fk_attempt_puzzle
    FOREIGN KEY (puzzle_id)
    REFERENCES puzzles(id)
    ON DELETE CASCADE
    
);
--Optimization Indexes
CREATE INDEX IF NOT EXISTS idx_attempts_user_id ON solve_attempts(user_id);
CREATE INDEX IF NOT EXISTS idx_attempts_puzzle_id ON solve_attempts(puzzle_id);
COMMIT;