-- +goose Up
-- +goose StatementBegin
BEGIN;

CREATE TABLE polls (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL
);

CREATE TABLE poll_votes (
    id SERIAL PRIMARY KEY,
    poll_id INTEGER NOT NULL REFERENCES polls(id) ON DELETE CASCADE,
    jti TEXT NOT NULL,
    vote TEXT NOT NULL CHECK (vote IN ('za', 'protiv'))
);

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

DROP TABLE IF EXISTS poll_votes;
DROP TABLE IF EXISTS polls;

COMMIT;
-- +goose StatementEnd
