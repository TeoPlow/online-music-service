from app.db.clickhouse import get_clickhouse_client

from app.core.logger import get_logger
log = get_logger(__name__)


async def run_clickhouse_migrations():
    """
    Запускает инициализацию БД в ClickHouse. Обычно при старте сервиса.
    """
    log.info("Running ClickHouse migrations...")
    client = await get_clickhouse_client()

    async with client.cursor() as cursor:

        try:
            log.info("Creating database 'music_streaming'...")
            await cursor.execute("""
                CREATE DATABASE IF NOT EXISTS music_streaming
            """)
            log.info("Database 'music_streaming' created or already exists.")

            # TRACKS
            await cursor.execute("""
                CREATE TABLE IF NOT EXISTS music_streaming.tracks (
                    id UUID,
                    title String,
                    album_id UUID,
                    genre String,
                    duration UInt32,
                    lyrics String,
                    is_explicit UInt8,
                    published_at DateTime,
                    updated_at DateTime
                )
                ENGINE = ReplacingMergeTree(updated_at)
                ORDER BY id
            """)
            log.info("Table 'tracks' created.")

            # ALBUMS
            await cursor.execute("""
                CREATE TABLE IF NOT EXISTS music_streaming.albums (
                    id UUID,
                    title String,
                    artist_id UUID,
                    release_date Date,
                    updated_at DateTime
                )
                ENGINE = ReplacingMergeTree(updated_at)
                ORDER BY id
            """)
            log.info("Table 'albums' created.")

            # ARTISTS
            await cursor.execute("""
                CREATE TABLE IF NOT EXISTS music_streaming.artists (
                    id UUID,
                    name String,
                    author String,
                    producer String,
                    country String,
                    description String,
                    created_at DateTime,
                    updated_at DateTime
                )
                ENGINE = ReplacingMergeTree(updated_at)
                ORDER BY id
            """)
            log.info("Table 'artists' created.")

            # USERS
            await cursor.execute("""
                CREATE TABLE IF NOT EXISTS music_streaming.users (
                    id UUID,
                    username String,
                    gender String,
                    country String,
                    age UInt16,
                    role String,
                    liked_tracks Array(UUID),
                    liked_artists Array(UUID),
                    created_at DateTime,
                    updated_at DateTime
                )
                ENGINE = ReplacingMergeTree(updated_at)
                ORDER BY id
            """)
            log.info("Table 'users' created.")

            # EVENTS
            await cursor.execute("""
                CREATE TABLE IF NOT EXISTS music_streaming.events (
                    event_time DateTime,
                    user_id UUID,
                    track_id UUID
                )
                ENGINE = MergeTree
                PARTITION BY toYYYYMM(event_time)
                ORDER BY (user_id, event_time)
            """)
            log.info("Table 'events' created.")
        except Exception as e:
            log.error(f"Failed to connect to database: {e}")

    log.info("ClickHouse migrations complete.")
