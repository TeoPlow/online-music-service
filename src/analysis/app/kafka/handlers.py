from uuid import UUID
from datetime import datetime

from asynch import DictCursor
from app.db.clickhouse import get_clickhouse_client
from app.db.models import User, Artist, Album, Track
from app.core.logger import get_logger

log = get_logger(__name__)


async def handle_auth_users(message: dict):
    try:
        client = await get_clickhouse_client()

        user_id = UUID(message["id"])
        old_user = await User.get_latest_by_id(user_id)

        liked_tracks = old_user.liked_tracks if old_user else []
        liked_artists = old_user.liked_artists if old_user else []

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
            liked_tracks=liked_tracks,
            liked_artists=liked_artists,
            created_at=message["created_at"],
            updated_at=datetime.now(),
        )

        async with client.cursor(cursor=DictCursor) as cursor:
            await cursor.execute("INSERT INTO music_streaming.users VALUES",
                                 [user.model_dump()])

    except Exception as e:
        log.error(f"Failed to handle auth-users message: {e}")


async def handler_music_liked_tracks(message: dict):
    try:
        client = await get_clickhouse_client()

        user_id = UUID(message["user_id"])
        new_track_id = UUID(message["track_id"])

        user = await User.get_latest_by_id(client, user_id)
        if not user:
            user = User(
                id=user_id,
                liked_tracks=[new_track_id],
                updated_at=datetime.now(),
            )
        else:
            if new_track_id not in user.liked_tracks:
                user.liked_tracks.append(new_track_id)
            user.updated_at = datetime.now()

        async with client.cursor(cursor=DictCursor) as cursor:
            await cursor.execute("INSERT INTO music_streaming.users VALUES",
                                 [user.model_dump()])

    except Exception as e:
        log.error(f"Failed to handle music-liked-tracks message: {e}")


async def handler_music_liked_artists(message: dict):
    try:
        client = await get_clickhouse_client()

        user_id = UUID(message["user_id"])
        new_artist_id = UUID(message["artist_id"])

        user = await User.get_latest_by_id(user_id)
        if not user:
            user = User(
                id=user_id,
                liked_artists=[new_artist_id],
                updated_at=datetime.now(),
            )
        else:
            if new_artist_id not in user.liked_artists:
                user.liked_artists.append(new_artist_id)
            user.updated_at = datetime.now()

        async with client.cursor(cursor=DictCursor) as cursor:
            await cursor.execute("INSERT INTO music_streaming.users VALUES",
                                 [user.model_dump()])

    except Exception as e:
        log.error(f"Failed to handle music-liked-artists message: {e}")


async def handle_music_artists(message: dict):
    try:
        client = await get_clickhouse_client()

        artist_id = UUID(message["id"])
        artist = Artist(
            id=artist_id,
            name=message["name"],
            author=message["author"],
            producer=message["producer"],
            country=message["country"],
            description=message["description"],
            created_at=message["created_at"],
            updated_at=datetime.now()
        )

        async with client.cursor(cursor=DictCursor) as cursor:
            await cursor.execute("INSERT INTO music_streaming.artists VALUES",
                                 [artist.model_dump()])

    except Exception as e:
        log.error(f"Failed to handle music-artists message: {e}")


async def handle_music_albums(message: dict):
    try:
        client = await get_clickhouse_client()

        album_id = UUID(message["id"])
        album = Album(
            id=album_id,
            title=message["title"],
            artist_id=UUID(message["artist_id"]),
            release_date=message["release_date"],
            updated_at=datetime.now()
        )

        async with client.cursor(cursor=DictCursor) as cursor:
            await cursor.execute("INSERT INTO music_streaming.albums VALUES",
                                 [album.model_dump()])

    except Exception as e:
        log.error(f"Failed to handle music-albums message: {e}")


async def handle_music_tracks(message: dict):
    try:
        client = await get_clickhouse_client()

        track_id = UUID(message["id"])
        track = Track(
            id=track_id,
            title=message["title"],
            album_id=UUID(message["album_id"]),
            genre=message["genre"],
            duration=message["duration"],
            lyrics=message["lyrics"],
            is_explicit=message["is_explicit"],
            published_at=message["published_at"],
            updated_at=datetime.now()
        )

        async with client.cursor(cursor=DictCursor) as cursor:
            await cursor.execute("INSERT INTO music_streaming.tracks VALUES",
                                 [track.model_dump()])

    except Exception as e:
        log.error(f"Failed to handle music-tracks message: {e}")
