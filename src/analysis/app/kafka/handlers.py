from uuid import UUID
from datetime import datetime

from app.db.connection import get_db_client
from app.db.models import (
    User,
    Artist,
    Album,
    Track,
    LikedArtist,
    LikedTrack,
    Event,
)
from app.core.logger import get_logger

log = get_logger(__name__)


async def handle_auth_users(message: dict):
    try:
        client = await get_db_client()

        user_id = UUID(message["id"])

        gender = "female"
        if message["gender"]:
            gender = "male"

        user = User(
            id=user_id,
            username=message["username"],
            gender=gender,
            country=message["country"],
            age=message["age"],
            role=message["role"],
            created_at=message["created_at"],
            updated_at=datetime.now(),
        )

        async with client.acquire() as conn:
            await conn.execute(
                """
                INSERT INTO music_streaming.users
                (id, username, gender, country,
                age, role, created_at, updated_at)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
                ON CONFLICT (id) DO UPDATE SET
                username = EXCLUDED.username,
                gender = EXCLUDED.gender,
                country = EXCLUDED.country,
                age = EXCLUDED.age,
                role = EXCLUDED.role,
                updated_at = EXCLUDED.updated_at
                """,
                user.id,
                user.username,
                user.gender,
                user.country,
                user.age,
                user.role,
                user.created_at,
                user.updated_at,
            )

    except Exception as e:
        log.error(f"Failed to handle auth-users message: {e}")


async def handler_music_liked_tracks(message: dict):
    try:
        client = await get_db_client()

        liked_track = LikedTrack(
            user_id=UUID(message["user_id"]),
            track_id=UUID(message["track_id"]),
        )

        async with client.acquire() as conn:
            await conn.execute(
                """
                INSERT INTO music_streaming.liked_tracks
                (user_id, track_id)
                VALUES ($1, $2)
                ON CONFLICT (user_id, track_id) DO NOTHING
                """,
                liked_track.user_id,
                liked_track.track_id,
            )

    except Exception as e:
        log.error(f"Failed to handle music-liked-tracks message: {e}")


async def handler_music_liked_artists(message: dict):
    try:
        client = await get_db_client()

        liked_artist = LikedArtist(
            user_id=UUID(message["user_id"]),
            artist_id=UUID(message["artist_id"]),
        )

        async with client.acquire() as conn:
            await conn.execute(
                """
                INSERT INTO music_streaming.liked_artists
                (user_id, artist_id)
                VALUES ($1, $2)
                ON CONFLICT (user_id, artist_id) DO NOTHING
                """,
                liked_artist.user_id,
                liked_artist.artist_id,
            )

    except Exception as e:
        log.error(f"Failed to handle music-liked-artists message: {e}")


async def handle_music_artists(message: dict):
    try:
        client = await get_db_client()

        artist_id = UUID(message["id"])
        artist = Artist(
            id=artist_id,
            name=message["name"],
            author=message["author"],
            producer=message["producer"],
            country=message["country"],
            description=message["description"],
            created_at=message["created_at"],
            updated_at=datetime.now(),
        )

        async with client.acquire() as conn:
            await conn.execute(
                """
                INSERT INTO music_streaming.artists
                (id, name, author, producer, country,
                description, created_at, updated_at)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
                ON CONFLICT (id) DO UPDATE SET
                name = EXCLUDED.name,
                author = EXCLUDED.author,
                producer = EXCLUDED.producer,
                country = EXCLUDED.country,
                description = EXCLUDED.description,
                updated_at = EXCLUDED.updated_at
                """,
                artist.id,
                artist.name,
                artist.author,
                artist.producer,
                artist.country,
                artist.description,
                artist.created_at,
                artist.updated_at,
            )

    except Exception as e:
        log.error(f"Failed to handle music-artists message: {e}")


async def handle_music_albums(message: dict):
    try:
        client = await get_db_client()

        album_id = UUID(message["id"])
        album = Album(
            id=album_id,
            title=message["title"],
            artist_id=UUID(message["artist_id"]),
            release_date=message["release_date"],
            updated_at=datetime.now(),
        )

        async with client.acquire() as conn:
            await conn.execute(
                """
                INSERT INTO music_streaming.albums
                (id, title, artist_id, release_date, updated_at)
                VALUES ($1, $2, $3, $4, $5)
                ON CONFLICT (id) DO UPDATE SET
                title = EXCLUDED.title,
                artist_id = EXCLUDED.artist_id,
                release_date = EXCLUDED.release_date,
                updated_at = EXCLUDED.updated_at
                """,
                album.id,
                album.title,
                album.artist_id,
                album.release_date,
                album.updated_at,
            )

    except Exception as e:
        log.error(f"Failed to handle music-albums message: {e}")


async def handle_music_tracks(message: dict):
    try:
        client = await get_db_client()

        track_id = UUID(message["id"])
        track = Track(
            id=track_id,
            title=message["title"],
            album_id=UUID(message["album_id"]),
            genre=message["genre"],
            duration=message["duration"],
            lyrics=message["lyrics"],
            is_explicit=message["is_explicit"],
            created_at=message["created_at"],
            updated_at=datetime.now(),
        )

        async with client.acquire() as conn:
            await conn.execute(
                """
                INSERT INTO music_streaming.tracks
                (id, title, album_id, genre, duration,
                lyrics, is_explicit, created_at, updated_at)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
                ON CONFLICT (id) DO UPDATE SET
                title = EXCLUDED.title,
                album_id = EXCLUDED.album_id,
                genre = EXCLUDED.genre,
                duration = EXCLUDED.duration,
                lyrics = EXCLUDED.lyrics,
                is_explicit = EXCLUDED.is_explicit,
                created_at = EXCLUDED.created_at,
                updated_at = EXCLUDED.updated_at
                """,
                track.id,
                track.title,
                track.album_id,
                track.genre,
                track.duration,
                track.lyrics,
                track.is_explicit,
                track.created_at,
                track.updated_at,
            )

    except Exception as e:
        log.error(f"Failed to handle music-tracks message: {e}")


async def handle_event(message: dict):
    try:
        client = await get_db_client()

        user_id = UUID(message["user_id"])
        track_id = UUID(message["track_id"])

        event = Event(
            event_time=datetime.now(),
            user_id=user_id,
            track_id=track_id,
        )

        async with client.acquire() as conn:
            await conn.execute(
                """
                INSERT INTO music_streaming.events
                (event_time, user_id, track_id)
                VALUES ($1, $2, $3)
                ON CONFLICT (user_id, track_id) DO NOTHING
                """,
                event.event_time,
                event.user_id,
                event.track_id,
            )

    except Exception as e:
        log.error(f"Failed to handle music-event message: {e}")
