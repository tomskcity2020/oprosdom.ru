-- +goose Up
-- +goose StatementBegin
BEGIN;

CREATE TABLE kvartiras (
    id         SERIAL       PRIMARY KEY,
    uuid       UUID         NOT NULL DEFAULT gen_random_uuid(),
    number     VARCHAR(10)  NOT NULL,
    komnat     INT          NOT NULL
);

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

DROP TABLE IF EXISTS kvartiras;

COMMIT;
-- +goose StatementEnd
