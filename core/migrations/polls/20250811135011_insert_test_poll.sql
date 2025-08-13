-- +goose Up
-- +goose StatementBegin
INSERT INTO polls (title) VALUES ('Установить видеонаблюдение 24 камеры');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM polls WHERE title = 'Установить видеонаблюдение 24 камеры';
-- +goose StatementEnd
