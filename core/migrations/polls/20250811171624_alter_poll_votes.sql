-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX idx_poll_votes_poll_jti ON poll_votes (poll_id, jti);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_poll_votes_poll_jti;
-- +goose StatementEnd
