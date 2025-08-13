-- +goose Up
-- +goose StatementBegin
INSERT INTO polls (title) VALUES ('Посадить газон вместо асфальта');
INSERT INTO polls (title) VALUES ('Сделать парковку вместо газона');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM polls WHERE title = 'Посадить газон вместо асфальта';
DELETE FROM polls WHERE title = 'Сделать парковку вместо газона';
-- +goose StatementEnd
