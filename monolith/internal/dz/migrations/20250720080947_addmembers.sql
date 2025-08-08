-- +goose Up
-- +goose StatementBegin
BEGIN;

CREATE TABLE members (
    id         SERIAL       PRIMARY KEY,
    uuid       UUID         NOT NULL DEFAULT gen_random_uuid(),
    fullname   VARCHAR(100) NOT NULL,
    phone      VARCHAR(20)  NOT NULL,
    community  INT          NOT NULL
);

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

DROP TABLE IF EXISTS members;

COMMIT;
-- +goose StatementEnd