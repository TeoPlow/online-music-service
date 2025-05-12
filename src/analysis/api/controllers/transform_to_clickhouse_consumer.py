import json
import threading
from kafka import KafkaConsumer
from datetime import datetime
from clickhouse_driver import Client

from configs.config import (
    kafka_bootstrap_servers,
    kafka_topics,
    clickhouse_host,
    clickhouse_port,
    clickhouse_user,
    clickhouse_password,
    clickhouse_database,
)

from api.models.clickhouse_models import (
    LikedArtist,
    LikedTrack,
    Song,
    Artist,
    User,
)

from utils.logger_config import configure
import logging
log = logging.getLogger('analysisLogger')
configure()

stop_event = threading.Event()

try:
    client = Client(
        clickhouse_host,
        port=clickhouse_port,
        user=clickhouse_user,
        password=clickhouse_password,
        database=clickhouse_database
    )

    # Костылёк для проверки работоспособности БД
    client.execute("SELECT 1")
    log.info("Successfully connected to ClickHouse")
except Exception as e:
    log.critical("Failed to connect to ClickHouse: %s", e)
    raise

album_cache = {}


def get_kafka_consumer() -> KafkaConsumer:
    """
    Создаёт и возвращает KafkaConsumer с настройками для обработки сообщений.
    """
    return KafkaConsumer(
        *kafka_topics,
        bootstrap_servers=kafka_bootstrap_servers,
        value_deserializer=lambda v: json.loads(v.decode("utf-8")),
        auto_offset_reset="earliest"
    )


def process_messages():
    """
    Основной цикл обработки сообщений из Kafka.
    """
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
            user = User(
                user_id=value["id"],
                gender="male" if value["gender"] else "female",
                age=value["age"],
                registration_date=datetime.fromisoformat(
                    value["created_at"]
                    ).date(),
                citizenship=value["country"]
            )
            user.save(client)

        elif topic == 'music-artists':
            artist = Artist(
                artist_id=value["id"],
                artist_name=value["name"],
                citizenship=value["country"]
            )
            artist.save(client)

        elif topic == 'music-albums':
            album_cache[value['id']] = value['artist_id']

        elif topic == 'music-tracks':
            artist_id = album_cache.get(value['album_id'])
            if artist_id:
                song = Song(
                    song_id=value["id"],
                    song_name=value["title"],
                    genre=value["genre"],
                    artist_id=artist_id
                )
                song.save(client)

        elif topic == 'music-liked-artists':
            liked_artist = LikedArtist(
                user_id=value["user_id"],
                artist_id=value["artist_id"]
            )
            liked_artist.save(client)

        elif topic == 'music-liked-tracks':
            liked_track = LikedTrack(
                user_id=value["user_id"],
                track_id=value["track_id"]
            )
            liked_track.save(client)


if __name__ == "__main__":
    process_messages()
