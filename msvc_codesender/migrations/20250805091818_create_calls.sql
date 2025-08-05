-- +goose Up
-- +goose StatementBegin
BEGIN;
CREATE TABLE IF NOT EXISTS call_messages (
    id SERIAL PRIMARY KEY,
    phone VARCHAR(20) NOT NULL, -- E.164 формат
    retry SMALLINT NOT NULL CHECK (retry BETWEEN 1 AND 3) DEFAULT 1,
    status VARCHAR(20) NOT NULL 
        CHECK (status IN ('wait', 'get', 'sent', 'expired')) 
        DEFAULT 'wait',
    gateway VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_call_messages_modtime
BEFORE UPDATE ON call_messages
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

CREATE TABLE IF NOT EXISTS calls_log (
    id SERIAL PRIMARY KEY,
    message_id INTEGER NOT NULL REFERENCES call_messages(id) ON DELETE CASCADE,
    worker_id VARCHAR(50) NOT NULL,
    gateway VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('attempt', 'sent', 'failed')),
    error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN call_messages.phone IS 'Номер в формате E.164';
COMMENT ON COLUMN call_messages.retry IS 'Допустимые значения: 1,2,3';
COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

-- Удаление триггера (должен быть удалён до функции)
DROP TRIGGER IF EXISTS update_call_messages_modtime ON call_messages;

-- Удаление дочерней таблицы (зависит от call_messages)
DROP TABLE IF EXISTS calls_log;

-- Удаление основной таблицы
DROP TABLE IF EXISTS call_messages;

COMMIT;
-- +goose StatementEnd
