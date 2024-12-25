-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders ADD updated_at TIMESTAMPTZ DEFAULT NOW();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP COLUMN updated_at;
-- +goose StatementEnd
