from uuid import UUID
from datetime import date, datetime
from typing import Optional
from pydantic import BaseModel

from app.core.logger import get_logger

log = get_logger(__name__)


class User(BaseModel):
    id: UUID
    username: str
    gender: Optional[str]
    country: Optional[str]
    age: Optional[int]
    role: str
    created_at: datetime
    updated_at: datetime

    @staticmethod
    async def get_latest_by_id(user_id: UUID, pool) -> Optional["User"]:
        log.info(f"Fetching latest user by ID: {user_id}")

        async with pool.acquire() as conn:
            query = """
                SELECT * FROM users
                WHERE id = $1
                ORDER BY updated_at DESC LIMIT 1
            """
            result = await conn.fetch(query, user_id)

        if not result:
            return None

        row = result[0]
        return User(
            id=row["id"],
            username=row["username"],
            gender=row["gender"],
            country=row["country"],
            age=row["age"],
            role=row["role"],
            created_at=row["created_at"],
            updated_at=row["updated_at"],
        )


class Track(BaseModel):
    id: UUID
    title: str
    album_id: UUID
    genre: Optional[str]
    duration: Optional[int]
    lyrics: Optional[str]
    is_explicit: bool
    created_at: datetime
    updated_at: datetime

    @staticmethod
    async def get_latest_by_id(track_id: UUID, pool) -> Optional["Track"]:
        log.info(f"Fetching latest track by ID: {track_id}")

        async with pool.acquire() as conn:
            query = """
                SELECT * FROM tracks
                WHERE id = $1
                ORDER BY updated_at DESC LIMIT 1
            """
            result = await conn.fetch(query, track_id)

        if not result:
            return None

        row = result[0]
        return Track(
            id=row["id"],
            title=row["title"],
            album_id=row["album_id"],
            genre=row["genre"],
            duration=row["duration"],
            lyrics=row["lyrics"],
            is_explicit=row["is_explicit"],
            created_at=row["created_at"],
            updated_at=row["updated_at"],
        )


class Album(BaseModel):
    id: UUID
    title: str
    artist_id: UUID
    release_date: date
    updated_at: datetime

    @staticmethod
    async def get_latest_by_id(album_id: UUID, pool) -> Optional["Album"]:
        log.info(f"Fetching latest album by ID: {album_id}")

        async with pool.acquire() as conn:
            query = """
                SELECT * FROM albums
                WHERE id = $1
                ORDER BY updated_at DESC LIMIT 1
            """
            result = await conn.fetch(query, album_id)

        if not result:
            return None

        row = result[0]
        return Album(
            id=row["id"],
            title=row["title"],
            artist_id=row["artist_id"],
            release_date=row["release_date"],
            updated_at=row["updated_at"],
        )


class Artist(BaseModel):
    id: UUID
    name: str
    author: str
    producer: str
    country: str
    description: str
    created_at: datetime
    updated_at: datetime

    @staticmethod
    async def get_latest_by_id(artist_id: UUID, pool) -> Optional["Artist"]:
        """
        Возвращает последнего добавленного в БД артиста с переданным UUID.
        """
        log.info(f"Fetching latest artist by ID: {artist_id}")

        async with pool.acquire() as conn:
            query = """
                SELECT * FROM artists
                WHERE id = $1
                ORDER BY updated_at DESC LIMIT 1
            """
            result = await conn.fetch(query, artist_id)

        if not result:
            return None

        row = result[0]
        return Artist(
            id=row["id"],
            name=row["name"],
            author=row["author"],
            producer=row["producer"],
            country=row["country"],
            description=row["description"],
            created_at=row["created_at"],
            updated_at=row["updated_at"],
        )


class Event(BaseModel):
    event_time: datetime
    user_id: UUID
    track_id: UUID

    @staticmethod
    async def get_latest_by_id(event_id: UUID, pool) -> Optional["Event"]:
        log.info(f"Fetching latest event by ID: {event_id}")

        async with pool.acquire() as conn:
            query = """
                SELECT * FROM events
                WHERE id = $1
                ORDER BY updated_at DESC LIMIT 1
            """
            result = await conn.fetch(query, event_id)

        if not result:
            return None

        row = result[0]
        return Event(
            event_time=row["event_time"],
            user_id=row["user_id"],
            track_id=row["track_id"],
        )


class LikedArtist(BaseModel):
    user_id: UUID
    artist_id: UUID

    @staticmethod
    async def get_latest_by_user_artist(
        user_id: UUID, artist_id: UUID, pool
    ) -> Optional["LikedArtist"]:
        async with pool.acquire() as conn:
            result = await conn.fetchrow(
                """
                SELECT * FROM liked_artists
                WHERE user_id = $1 AND artist_id = $2
                LIMIT 1
            """,
                user_id,
                artist_id,
            )
        if not result:
            return None
        return LikedArtist(**result)


class LikedTrack(BaseModel):
    user_id: UUID
    track_id: UUID

    @staticmethod
    async def get_latest_by_user_track(
        user_id: UUID, track_id: UUID, pool
    ) -> Optional["LikedTrack"]:
        async with pool.acquire() as conn:
            result = await conn.fetchrow(
                """
                SELECT * FROM liked_tracks
                WHERE user_id = $1 AND track_id = $2
                LIMIT 1
            """,
                user_id,
                track_id,
            )
        if not result:
            return None
        return LikedTrack(**result)
