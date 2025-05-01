-- +goose Up
-- +goose StatementBegin
CREATE TABLE tracks (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    album_id UUID NOT NULL,
    genre TEXT NOT NULL,
    duration INTEGER NOT NULL CHECK (duration >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tracks;
-- +goose StatementEnd
