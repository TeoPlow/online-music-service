import random
import uuid
from datetime import datetime
from aiokafka import AIOKafkaProducer
import json
from faker import Faker
import asyncio
import os

from app.kafka.topics import (
    MUSIC_ARTISTS,
    MUSIC_ALBUMS,
    MUSIC_TRACKS,
    MUSIC_LIKED_TRACKS,
    MUSIC_LIKED_ARTISTS,
    AUTH_USERS,
)

from app.core.settings import settings

from app.core.logger import get_logger

log = get_logger(__name__)
fake = Faker()

log.info("[musical] Started async fake musical producer.")

try:
    with open("../static/genres.json", "r") as f:
        genres = json.load(f)
except FileNotFoundError:
    log.error(f"File {os.getcwd()}/../static/genres.json not found.")
    log.error("Exiting the program..")
    exit(1)


async def send_musical_entities(num_iterations=1):
    """
    Это функция с запуском фейкового асинхронного Kafka Producer,
    который будет спамить заданное кол-во итераций автосгенерированных
    пробросов Артистов, Треков, Альбомов, Понравившихся Треков и Артистов.
    """
    producer = AIOKafkaProducer(
        bootstrap_servers=settings.kafka.bootstrap_servers,
        value_serializer=lambda v: json.dumps(v).encode("utf-8"),
    )
    await producer.start()

    try:
        # Создаю одного тестового пользователя.
        with open(
            f"{settings.paths.base_dir}/../static/countries.json", "r"
        ) as f:
            countries = json.load(f)
        roles = ["user", "admin", "artist"]

        user_id = str(uuid.uuid4())
        user = {
            "id": user_id,
            "username": fake.user_name(),
            "email": fake.email(),
            "gender": random.choice([True, False]),
            "country": random.choice(countries),
            "age": random.randint(14, 80),
            "role": random.choice(roles),
            "passHash": fake.sha256(),
            "created_at": datetime.now().isoformat().split(".")[0],
        }
        await producer.send(AUTH_USERS, user)

        for _ in range(num_iterations):
            artists = []
            albums = []
            tracks = []
            liked_artists = []
            liked_tracks = []

            artist_id = str(uuid.uuid4())
            artist = {
                "id": artist_id,
                "name": fake.name(),
                "author": fake.name(),
                "producer": fake.name(),
                "country": fake.country_code(),
                "description": fake.sentence(),
                "created_at": datetime.now().isoformat().split(".")[0],
                "updated_at": datetime.now().isoformat().split(".")[0],
            }
            artists.append(artist)

            album_id = str(uuid.uuid4())
            album = {
                "id": album_id,
                "title": fake.catch_phrase(),
                "artist_id": artist_id,
                "release_date": fake.date_this_decade()
                .isoformat()
                .split(".")[0],
            }
            albums.append(album)

            track = {
                "id": str(uuid.uuid4()),
                "title": fake.sentence(nb_words=3),
                "album_id": album_id,
                "genre": random.choice(genres),
                "duration": random.randint(120, 420),
                "lyrics": (
                    fake.text(max_nb_chars=100)
                    if random.choice([True, False])
                    else None
                ),
                "is_explicit": random.choice([True, False]),
                "created_at": datetime.now().isoformat().split(".")[0],
                "updated_at": datetime.now().isoformat().split(".")[0],
            }
            tracks.append(track)

            liked_artist = {
                "user_id": user_id,
                "artist_id": artist_id,
                "created_at": datetime.now().isoformat().split(".")[0],
            }
            liked_artists.append(liked_artist)

            liked_track = {
                "user_id": user_id,
                "track_id": track["id"],
                "created_at": datetime.now().isoformat().split(".")[0],
            }
            liked_tracks.append(liked_track)

            # Отправка данных в Kafka
            for artist in artists:
                await producer.send(MUSIC_ARTISTS, artist)
                log.debug(f"[{MUSIC_ARTISTS}] Kafka Sent: {artist}")
                await asyncio.sleep(0.5)

            for album in albums:
                await producer.send(MUSIC_ALBUMS, album)
                log.debug(f"[{MUSIC_ALBUMS}] Kafka Sent: {album}")
                await asyncio.sleep(0.5)

            for track in tracks:
                await producer.send(MUSIC_TRACKS, track)
                log.debug(f"[{MUSIC_TRACKS}] Kafka Sent: {track}")
                await asyncio.sleep(0.5)

            for liked_artist in liked_artists:
                await producer.send(MUSIC_LIKED_ARTISTS, liked_artist)
                log.debug(f"[{MUSIC_LIKED_ARTISTS}] Kafka Sent:{liked_artist}")
                await asyncio.sleep(0.5)

            for liked_track in liked_tracks:
                await producer.send(MUSIC_LIKED_TRACKS, liked_track)
                log.debug(f"[{MUSIC_LIKED_TRACKS}] Kafka Sent: {liked_track}")
                await asyncio.sleep(0.5)

    finally:
        await producer.stop()
    return artists, albums, tracks, liked_artists, liked_tracks
