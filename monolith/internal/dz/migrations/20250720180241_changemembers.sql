-- +goose Up
-- +goose StatementBegin
ALTER TABLE members RENAME COLUMN fullname TO name;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE members RENAME COLUMN name TO fullname;
-- +goose StatementEnd
