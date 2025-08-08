-- +goose Up
-- +goose StatementBegin
BEGIN;
CREATE TABLE IF NOT EXISTS sms_messages (
    id SERIAL PRIMARY KEY,
    phone VARCHAR(20) NOT NULL, -- E.164 формат
    message TEXT NOT NULL,
    retry SMALLINT NOT NULL CHECK (retry BETWEEN 1 AND 3) DEFAULT 1,
    status VARCHAR(20) NOT NULL 
        CHECK (status IN ('wait', 'get', 'sent', 'expired')) 
        DEFAULT 'wait',
    gateway VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_sms_messages_modtime
BEFORE UPDATE ON sms_messages
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

CREATE TABLE IF NOT EXISTS sms_log (
    id SERIAL PRIMARY KEY,
    message_id INTEGER NOT NULL REFERENCES sms_messages(id) ON DELETE CASCADE,
    worker_id VARCHAR(50) NOT NULL,
    gateway VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('attempt', 'sent', 'failed')),
    error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN sms_messages.phone IS 'Номер в формате E.164';
COMMENT ON COLUMN sms_messages.retry IS 'Допустимые значения: 1,2,3';
COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

-- Удаление триггера (должен быть удалён до функции)
DROP TRIGGER IF EXISTS update_sms_messages_modtime ON sms_messages;

-- Удаление дочерней таблицы (зависит от sms_messages)
DROP TABLE IF EXISTS sms_log;

-- Удаление основной таблицы
DROP TABLE IF EXISTS sms_messages;

-- Удаление функции триггера (после удаления триггера)
DROP FUNCTION IF EXISTS update_modified_column;

COMMIT;
-- +goose StatementEnd
