-- +goose Up
-- +goose StatementBegin
BEGIN;

-- Добавляем NOT NULL столбец code в sms_messages
ALTER TABLE sms_messages 
ADD COLUMN code INTEGER NOT NULL;

COMMENT ON COLUMN sms_messages.code IS 'Код подтверждения от 1000 до 9999';

-- Добавляем NOT NULL столбец code в call_messages
ALTER TABLE call_messages 
ADD COLUMN code INTEGER NOT NULL;

COMMENT ON COLUMN call_messages.code IS 'Код подтверждения от 1000 до 9999';

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

-- Удаляем столбец code из sms_messages
ALTER TABLE sms_messages 
DROP COLUMN IF EXISTS code;

-- Удаляем столбец code из call_messages
ALTER TABLE call_messages 
DROP COLUMN IF EXISTS code;

COMMIT;
-- +goose StatementEnd