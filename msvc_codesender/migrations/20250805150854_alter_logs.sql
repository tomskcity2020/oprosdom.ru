-- +goose Up
-- +goose StatementBegin
BEGIN;

ALTER TABLE sms_log DROP COLUMN worker_id;
ALTER TABLE calls_log DROP COLUMN worker_id;

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

ALTER TABLE sms_log ADD COLUMN worker_id TEXT NOT NULL DEFAULT '';
ALTER TABLE calls_log ADD COLUMN worker_id TEXT NOT NULL DEFAULT '';

COMMIT;
-- +goose StatementEnd
