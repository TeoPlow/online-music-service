-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY,
            username TEXT,
            gender TEXT,
            country TEXT,
            age SMALLINT,
            role TEXT,
            created_at TIMESTAMP,
            updated_at TIMESTAMP
        );

CREATE TABLE IF NOT EXISTS artists (
            id UUID PRIMARY KEY,
            name TEXT,
            author TEXT,
            producer TEXT,
            country TEXT,
            description TEXT,
            created_at TIMESTAMP,
            updated_at TIMESTAMP
        );

CREATE TABLE IF NOT EXISTS albums (
            id UUID PRIMARY KEY,
            title TEXT,
            artist_id UUID REFERENCES artists(id),
            release_date DATE,
            updated_at TIMESTAMP
        );

CREATE TABLE IF NOT EXISTS tracks (
            id UUID PRIMARY KEY,
            title TEXT,
            album_id UUID REFERENCES albums(id),
            genre TEXT,
            duration INTEGER,
            lyrics TEXT,
            is_explicit BOOLEAN,
            created_at TIMESTAMP,
            updated_at TIMESTAMP
        );

CREATE TABLE IF NOT EXISTS events (
            event_time TIMESTAMP,
            user_id UUID REFERENCES users(id),
            track_id UUID REFERENCES tracks(id)
        );

CREATE TABLE IF NOT EXISTS liked_tracks (
            user_id UUID REFERENCES users(id),
            track_id UUID REFERENCES tracks(id),
            PRIMARY KEY (user_id, track_id)
        );

CREATE TABLE IF NOT EXISTS liked_artists (
            user_id UUID REFERENCES users(id),
            artist_id UUID REFERENCES artists(id),
            PRIMARY KEY (user_id, artist_id)
        );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS liked_artists;
DROP TABLE IF EXISTS liked_tracks;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS tracks;
DROP TABLE IF EXISTS albums;
DROP TABLE IF EXISTS artists;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd