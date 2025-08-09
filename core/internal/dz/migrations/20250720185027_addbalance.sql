-- +goose Up
-- +goose StatementBegin
ALTER TABLE members 
ADD COLUMN balance DECIMAL(15, 2) NOT NULL DEFAULT 100000.00;
ALTER TABLE kvartiras 
ADD COLUMN debt DECIMAL(15, 2) NOT NULL DEFAULT 120000.00;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE members DROP COLUMN balance;
ALTER TABLE kvartiras DROP COLUMN debt;
-- +goose StatementEnd
