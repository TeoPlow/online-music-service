CREATE DATABASE IF NOT EXISTS music_streaming;

-- Users
CREATE TABLE IF NOT EXISTS music_streaming.users (
    user_id UUID,
    gender String,
    age UInt8,
    registration_date Date,
    citizenship String,
    liked_songs Array(UUID),
    disliked_songs Array(UUID),
    liked_artists Array(UUID),
    disliked_artists Array(UUID)
) ENGINE = MergeTree ORDER BY user_id;

-- Songs
CREATE TABLE IF NOT EXISTS music_streaming.songs (
    song_id UUID,
    song_name String,
    genre String,
    artist_id UUID
) ENGINE = MergeTree ORDER BY song_id;

-- Artists
CREATE TABLE IF NOT EXISTS music_streaming.artists (
    artist_id UUID,
    artist_name String,
    citizenship String
) ENGINE = MergeTree ORDER BY artist_id;

-- Events
CREATE TABLE IF NOT EXISTS music_streaming.events (
    event_time DateTime,
    user_id UUID,
    song_id UUID
) ENGINE = MergeTree
PARTITION BY toYYYYMM(event_time)
ORDER BY (user_id, event_time);
