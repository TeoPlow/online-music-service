from uuid import UUID
from datetime import date, datetime
from typing import Optional
from pydantic import BaseModel

from app.core.logger import get_logger

log = get_logger(__name__)


class User(BaseModel):
    id: UUID
    username: str
    gender: str
    country: str
    age: int
    role: str
    created_at: datetime
    updated_at: datetime

    @staticmethod
    async def get_latest_by_id(user_id: UUID, pool) -> Optional["User"]:
        log.info(f"Fetching latest user by ID: {user_id}")

        async with pool.acquire() as conn:
            query = """
                SELECT * FROM music_streaming.users
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
    genre: str
    duration: int
    lyrics: str
    is_explicit: bool
    published_at: datetime
    updated_at: datetime


class Album(BaseModel):
    id: UUID
    title: str
    artist_id: UUID
    release_date: date
    updated_at: datetime


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
                SELECT * FROM music_streaming.artists
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


class LikedTrack(BaseModel):
    user_id: UUID
    track_id: UUID


class LikedArtist(BaseModel):
    user_id: UUID
    artist_id: UUID
