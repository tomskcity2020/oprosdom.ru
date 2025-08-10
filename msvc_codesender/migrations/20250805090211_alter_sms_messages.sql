-- +goose Up
-- +goose StatementBegin
BEGIN;

ALTER TABLE sms_messages DROP COLUMN message;

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

ALTER TABLE sms_messages ADD COLUMN message TEXT NOT NULL DEFAULT '';

COMMIT;
-- +goose StatementEnd
