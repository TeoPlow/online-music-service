import asyncio
import json
from aiokafka import AIOKafkaConsumer

from app.kafka.handlers import (
    handle_auth_users,
    handle_music_artists,
    handle_music_albums,
    handle_music_tracks,
    handler_music_liked_tracks,
    handler_music_liked_artists,
)
from app.kafka.topics import (
    AUTH_USERS,
    MUSIC_ARTISTS,
    MUSIC_ALBUMS,
    MUSIC_TRACKS,
    MUSIC_LIKED_TRACKS,
    MUSIC_LIKED_ARTISTS,
)

from app.core.settings import settings

from app.core.logger import get_logger

log = get_logger(__name__)

TOPIC_HANDLERS = {
    AUTH_USERS: handle_auth_users,
    MUSIC_ARTISTS: handle_music_artists,
    MUSIC_ALBUMS: handle_music_albums,
    MUSIC_TRACKS: handle_music_tracks,
    MUSIC_LIKED_TRACKS: handler_music_liked_tracks,
    MUSIC_LIKED_ARTISTS: handler_music_liked_artists,
}


async def consume():
    log.info("Starting Kafka Consumer...")
    consumer = AIOKafkaConsumer(
        *settings.kafka.topics,
        bootstrap_servers=settings.kafka.bootstrap_servers,
        # group_id="analysis-service-group",
        value_deserializer=lambda v: json.loads(v.decode("utf-8")),
        # enable_auto_commit=True,
        auto_offset_reset="earliest",
    )

    await consumer.start()
    log.info(f"Kafka consumer started on {settings.kafka.bootstrap_servers}")
    try:
        async for msg in consumer:
            topic = msg.topic
            value = msg.value

            handler = TOPIC_HANDLERS.get(topic)
            if handler:
                log.info(f"Received message from topic: {topic}")
                try:
                    await handler(value)
                except Exception as e:
                    log.error(f"Handler for topic {topic} failed: {e}")
            else:
                log.warning(f"No handler defined for topic: {topic}")
    except Exception as e:
        log.error(f"Error in Kafka consumer: {e}")
    finally:
        await consumer.stop()
        log.info("Kafka consumer stopped.")


if __name__ == "__main__":
    asyncio.run(consume())
