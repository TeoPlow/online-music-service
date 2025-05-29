CREATE SCHEMA IF NOT EXISTS music_streaming;

CREATE TABLE IF NOT EXISTS music_streaming.users (
            id UUID PRIMARY KEY,
            username TEXT,
            gender TEXT,
            country TEXT,
            age SMALLINT,
            role TEXT,
            created_at TIMESTAMP,
            updated_at TIMESTAMP
        );

CREATE TABLE IF NOT EXISTS music_streaming.artists (
            id UUID PRIMARY KEY,
            name TEXT,
            author TEXT,
            producer TEXT,
            country TEXT,
            description TEXT,
            created_at TIMESTAMP,
            updated_at TIMESTAMP
        );

CREATE TABLE IF NOT EXISTS music_streaming.albums (
            id UUID PRIMARY KEY,
            title TEXT,
            artist_id UUID REFERENCES music_streaming.artists(id),
            release_date DATE,
            updated_at TIMESTAMP
        );

CREATE TABLE IF NOT EXISTS music_streaming.tracks (
            id UUID PRIMARY KEY,
            title TEXT,
            album_id UUID REFERENCES music_streaming.albums(id),
            genre TEXT,
            duration INTEGER,
            lyrics TEXT,
            is_explicit BOOLEAN,
            published_at TIMESTAMP,
            updated_at TIMESTAMP
        );

CREATE TABLE IF NOT EXISTS music_streaming.events (
            event_time TIMESTAMP,
            user_id UUID REFERENCES music_streaming.users(id),
            track_id UUID REFERENCES music_streaming.tracks(id)
        );

CREATE TABLE IF NOT EXISTS music_streaming.liked_tracks (
            user_id UUID REFERENCES music_streaming.users(id),
            track_id UUID REFERENCES music_streaming.tracks(id),
            PRIMARY KEY (user_id, track_id)
        );

CREATE TABLE IF NOT EXISTS music_streaming.liked_artists (
            user_id UUID REFERENCES music_streaming.users(id),
            artist_id UUID REFERENCES music_streaming.artists(id),
            PRIMARY KEY (user_id, artist_id)
        );
