-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS courier_latest_cord (
               courier_id UUID NOT NULL,
               latitude double precision NOT NULL,
               longitude double precision NOT NULL ,
               created_at TIMESTAMPTZ NOT NULL,
               PRIMARY KEY (courier_id, created_at)
    )

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP DATABASE courier_latest_cord;
-- +goose StatementEnd
