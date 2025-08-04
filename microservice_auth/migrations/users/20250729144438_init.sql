-- +goose Up
-- +goose StatementBegin
CREATE TABLE phonesend (
    id SERIAL PRIMARY KEY,
    phone VARCHAR(20) NOT NULL,
    useragent VARCHAR(550) NOT NULL DEFAULT '',
    ip INET NOT NULL,
    time BIGINT NOT NULL DEFAULT (EXTRACT(EPOCH FROM now())::BIGINT)
);

CREATE INDEX idx_phonesend_ip ON phonesend (ip);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_phonesend_ip;
DROP TABLE IF EXISTS phonesend;
-- +goose StatementEnd
