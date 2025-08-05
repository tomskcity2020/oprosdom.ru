-- +goose Up
-- +goose StatementBegin
BEGIN;
-- Создаем тип ENUM
CREATE TYPE phone_type AS ENUM ('mobile', 'landline', 'unknown');

-- Добавляем столбец с ENUM
ALTER TABLE phonesend
ADD COLUMN phone_type phone_type NOT NULL DEFAULT 'unknown';
COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;
-- Удаляем столбец
ALTER TABLE phonesend
DROP COLUMN phone_type;

-- Удаляем тип ENUM
DROP TYPE phone_type;
COMMIT;
-- +goose StatementEnd
