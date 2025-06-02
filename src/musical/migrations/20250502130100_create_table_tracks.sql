-- +goose Up
-- +goose StatementBegin
CREATE TABLE "tracks"(
    "id" UUID PRIMARY KEY,
    "title" VARCHAR(255) NOT NULL,
    "album_id" UUID NOT NULL,
    "genre" VARCHAR(255) NOT NULL,
    "duration" BIGINT NOT NULL,
    "lyrics" TEXT NULL,
    "is_explicit" BOOLEAN NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL
);

CREATE TABLE "albums"(
    "id" UUID PRIMARY KEY,
    "title" VARCHAR(255) NOT NULL,
    "artist_id" UUID NOT NULL,
    "release_date" DATE NOT NULL
);

CREATE TABLE "artists"(
    "id" UUID PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "author" VARCHAR(255) NULL,
    "producer" VARCHAR(255) NULL,
    "country" VARCHAR(255) NOT NULL,
    "description" VARCHAR(255) NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL
);

CREATE TABLE "liked_artists"(
    "user_id" UUID NOT NULL,
    "artist_id" UUID NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    CONSTRAINT fk_artist FOREIGN KEY (artist_id) REFERENCES artists(id) ON DELETE CASCADE

);
CREATE TABLE "liked_tracks"(
    "user_id" UUID NOT NULL,
    "track_id" UUID NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    CONSTRAINT fk_track FOREIGN KEY (track_id) REFERENCES tracks(id) ON DELETE CASCADE

);

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX trgm_idx_tracks_title ON tracks USING gin (title gin_trgm_ops);
CREATE INDEX trgm_idx_tracks_lyrics ON tracks USING gin (lyrics gin_trgm_ops);

CREATE INDEX trgm_idx_albums_title ON albums USING gin (title gin_trgm_ops);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "tracks";
DROP TABLE IF EXISTS "albums";
DROP TABLE IF EXISTS "artists";
DROP TABLE IF EXISTS "liked_artists";
DROP TABLE IF EXISTS "liked_tracks";
-- +goose StatementEnd
