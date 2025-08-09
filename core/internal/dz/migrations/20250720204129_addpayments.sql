-- +goose Up
-- +goose StatementBegin
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    member_uuid UUID NOT NULL,
    kvartira_uuid UUID NOT NULL,
    amount DECIMAL(15,2) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payment;
-- +goose StatementEnd
