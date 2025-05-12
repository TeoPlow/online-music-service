import time
import random
import uuid
from datetime import datetime
from kafka import KafkaProducer
import json
from faker import Faker
import os

from configs.config import (
    kafka_bootstrap_servers,
)

from utils.logger_config import configure
import logging
log = logging.getLogger('analysisLogger')
configure()

fake = Faker()

log.info("[musical] Запущен фейковый продюсер.")

producer = KafkaProducer(
    bootstrap_servers=kafka_bootstrap_servers,
    value_serializer=lambda v: json.dumps(v).encode('utf-8')
)

try:
    with open('../static/genres.json', 'r') as f:
        genres = json.load(f)

except FileNotFoundError:
    log.error(f"Файл {os.getcwd()}/../static/genres.json не найден.")
    exit(1)

artists = []
albums = []
tracks = []
liked_artists = []
liked_tracks = []


def send_music_entities(num_iterations=1):
    for _ in range(num_iterations):
        artists = []
        albums = []
        tracks = []
        liked_artists = []
        liked_tracks = []

        for _ in range(5):
            artist_id = str(uuid.uuid4())
            artist = {
                "id": artist_id,
                "name": fake.name(),
                "author": fake.name(),
                "producer": fake.name(),
                "country": fake.country_code(),
                "description": fake.sentence(),
                "created_at": datetime.now().isoformat(),
                "updated_at": datetime.now().isoformat()
            }
            artists.append(artist)

            album_id = str(uuid.uuid4())
            album = {
                "id": album_id,
                "title": fake.catch_phrase(),
                "artist_id": artist_id,
                "release_date": fake.date_this_decade().isoformat()
            }
            albums.append(album)

            track = {
                "id": str(uuid.uuid4()),
                "title": fake.sentence(nb_words=3),
                "album_id": album_id,
                "genre": random.choice(genres),
                "duration": random.randint(120, 420),
                "lyrics": fake.text(
                    max_nb_chars=100
                ) if random.choice([True, False]) else None,
                "is_explicit": random.choice([True, False]),
                "created_at": datetime.now().isoformat(),
                "updated_at": datetime.now().isoformat()
            }
            tracks.append(track)

            liked_artist = {
                "user_id": str(uuid.uuid4()),
                "artist_id": artist_id,
                "created_at": datetime.now().isoformat()
            }
            liked_artists.append(liked_artist)

            liked_track = {
                "user_id": str(uuid.uuid4()),
                "track_id": track["id"],
                "created_at": datetime.now().isoformat()
            }
            liked_tracks.append(liked_track)

        for artist in artists:
            producer.send('music-artists', artist)
            log.info(f"[music-artist] Kafka Sent: {artist}")
            time.sleep(1)

        for album in albums:
            producer.send('music-albums', album)
            log.info(f"[music-album] Kafka Sent: {album}")
            time.sleep(1)

        for track in tracks:
            producer.send('music-tracks', track)
            log.info(f"[music-track] Kafka Sent: {track}")
            time.sleep(random.randint(5, 10))

        for liked_artist in liked_artists:
            producer.send('music-liked-artists', liked_artist)
            log.info(f"[music-liked-artist] Kafka Sent: {liked_artist}")
            time.sleep(1)

        for liked_track in liked_tracks:
            producer.send('music-liked-tracks', liked_track)
            log.info(f"[music-liked-track] Kafka Sent: {liked_track}")
            time.sleep(1)


if __name__ == '__main__':
    while True:
        send_music_entities()
