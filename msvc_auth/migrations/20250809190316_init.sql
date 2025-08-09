-- +goose Up
-- +goose StatementBegin
BEGIN;

CREATE TABLE signed_tokens (
    jti UUID PRIMARY KEY,
    alg VARCHAR(10) NOT NULL,
    pubkey_id CHAR(8) NOT NULL,
    ident VARCHAR(50) NOT NULL,
    value TEXT NOT NULL,
    ip INET NOT NULL,
    user_agent VARCHAR(515) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_signed_tokens_pubkey_id ON signed_tokens(pubkey_id);
CREATE INDEX idx_signed_tokens_ident ON signed_tokens(ident);
CREATE INDEX idx_signed_tokens_ip ON signed_tokens(ip);

COMMIT;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

DROP TABLE IF EXISTS signed_tokens;

COMMIT;
-- +goose StatementEnd
