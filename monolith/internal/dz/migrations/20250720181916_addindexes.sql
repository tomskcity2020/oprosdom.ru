-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_members_community ON members(community);
CREATE INDEX idx_kvartiras_number ON kvartiras(number);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_members_community;
DROP INDEX IF EXISTS idx_kvartiras_number;
-- +goose StatementEnd
