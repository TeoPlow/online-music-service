import json
from kafka import KafkaConsumer
from datetime import datetime
from clickhouse_driver import Client

from src.analysis.configs.config import (
    kafka_bootstrap_servers,
    kafka_topics,
    clickhouse_host,
    clickhouse_port,
    clickhouse_user,
    clickhouse_password,
    clickhouse_database
)

from src.analysis.utils.logger_config import configure
import logging
log = logging.getLogger('analysisLogger')
configure()

try:
    client = Client(
        clickhouse_host,
        port=clickhouse_port,
        user=clickhouse_user,
        password=clickhouse_password,
        database=clickhouse_database
    )
    client.execute("SELECT 1")
    log.info("Successfully connected to ClickHouse")
except Exception as e:
    log.critical("Failed to connect to ClickHouse: %s", e)
    raise

album_cache = {}


def get_kafka_consumer():
    return KafkaConsumer(
        *kafka_topics,
        bootstrap_servers=kafka_bootstrap_servers,
        value_deserializer=lambda v: json.loads(v.decode("utf-8")),
        auto_offset_reset="earliest"
    )


def insert_user(user):
    try:
        client.execute(
            """
            INSERT INTO music_streaming.users (
                user_id,
                gender,
                age,
                registration_date,
                citizenship,
                liked_songs,
                liked_artists
            ) VALUES
            """,
            [(
                user["id"],
                "male" if user["gender"] else "female",
                user["age"],
                datetime.fromisoformat(user["created_at"]).date(),
                user["country"],
                [],
                []
            )]
        )
        log.info("User inserted successfully: %s", user["id"])
    except Exception as e:
        log.error("Failed to insert user: %s", e)


def insert_artist(artist):
    try:
        client.execute(
            """
            INSERT INTO music_streaming.artists (
                artist_id,
                artist_name,
                citizenship
            ) VALUES
            """,
            [(
                artist["id"],
                artist["name"],
                artist["country"]
            )]
        )
        log.info("Artist inserted successfully: %s", artist["id"])
    except Exception as e:
        log.error("Failed to insert artist: %s", e)


def insert_song(track, artist_id):
    try:
        client.execute(
            """
            INSERT INTO music_streaming.songs (
                song_id,
                song_name,
                genre,
                artist_id
            ) VALUES
            """,
            [(
                track["id"],
                track["title"],
                track["genre"],
                artist_id
            )]
        )
        log.info("Song inserted successfully: %s", track["id"])
    except Exception as e:
        log.error("Failed to insert song: %s", e)


def insert_liked_artist(liked_artist):
    try:
        client.execute(
            """
            INSERT INTO music_streaming.users (
                user_id,
                liked_artists
            ) VALUES
            """,
            [(
                liked_artist["user_id"],
                [liked_artist["artist_id"]]
            )]
        )
        log.info("Liked artist added success: %s", liked_artist["artist_id"])
    except Exception as e:
        log.error("Failed to add liked artist: %s", e)


def insert_liked_track(liked_track):
    try:
        client.execute(
            """
            INSERT INTO music_streaming.users (
                user_id,
                liked_songs
            ) VALUES
            """,
            [(
                liked_track["user_id"],
                [liked_track["track_id"]]
            )]
        )
        log.info("Liked track added successfully: %s", liked_track["track_id"])
    except Exception as e:
        log.error("Failed to add liked track: %s", e)


def process_messages():
    consumer = get_kafka_consumer()
    for message in consumer:
        try:
            topic = message.topic
            value = message.value
            log.info("Received message from topic: %s, message: %s",
                     topic,
                     value
                     )
        except Exception as e:
            log.error("Failed to process message: %s", e)

        if topic == 'auth-users':
            insert_user(value)

        elif topic == 'music-artists':
            insert_artist(value)

        elif topic == 'music-albums':
            album_cache[value['id']] = value['artist_id']

        elif topic == 'music-tracks':
            artist_id = album_cache.get(value['album_id'])
            if artist_id:
                insert_song(value, artist_id)

        elif topic == 'music-liked-artists':
            insert_liked_artist(value)

        elif topic == 'music-liked-tracks':
            insert_liked_track(value)


if __name__ == "__main__":
    process_messages()
