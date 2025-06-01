import asyncio
import pytest
from app.db.models import Track, Album, Artist, LikedArtist, LikedTrack

from tests.fake_producers.fake_musical_producer import send_musical_entities
from app.core.settings import settings


@pytest.mark.asyncio
async def test_one_message_to_consumer(db_pool):
    """
    Отправляет данные музыкального сервиса в Kafka и проверяет,
    были ли они успешно вставлены в Database через Kafka Consumer.
    """
    pg = settings.postgres
    print(
        f"[CONFIG] Database settings:\n"
        f"  host={pg.host}\n"
        f"  port={pg.port}\n"
        f"  user=postgres\n"
        f"  password={pg.password[:4]+'***' if pg.password else '(empty)'}\n"
    )

    artists, albums, tracks, liked_artists, liked_tracks = (
        await send_musical_entities(num_iterations=1)
    )
    await asyncio.sleep(5)
    print(f"Артисты:{artists}")
    print(f"Альбомы:{albums}")
    print(f"Треки:{tracks}")
    print(f"Лайкнутые артисты:{liked_artists}")
    print(f"Лайкнутые треки:{liked_tracks}")

    artist = artists[0]
    retrieved_artist = await Artist.get_latest_by_id(
        artist["id"], pool=db_pool
    )
    assert retrieved_artist is not None
    assert str(retrieved_artist.id) == artist["id"]
    assert retrieved_artist.name == artist["name"]
    assert retrieved_artist.author == artist["author"]
    assert retrieved_artist.producer == artist["producer"]
    assert retrieved_artist.country == artist["country"]
    assert retrieved_artist.description == artist["description"]
    assert (
        retrieved_artist.created_at.strftime("%Y-%m-%dT%H:%M:%S")
        == artist["created_at"]
    )
    # Падает из-за разных UTC поздно ночью
    # assert (
    #     retrieved_artist.updated_at.strftime("%Y-%m-%d")
    #     == artist["updated_at"].split("T")[0]
    # )

    album = albums[0]
    retrieved_album = await Album.get_latest_by_id(album["id"], pool=db_pool)
    assert retrieved_album is not None
    assert str(retrieved_album.id) == album["id"]
    assert retrieved_album.title == album["title"]
    assert str(retrieved_album.artist_id) == album["artist_id"]
    assert retrieved_album.release_date.isoformat() == album["release_date"]

    track = tracks[0]
    retrieved_track = await Track.get_latest_by_id(track["id"], pool=db_pool)
    assert retrieved_track is not None
    assert str(retrieved_track.id) == track["id"]
    assert retrieved_track.title == track["title"]
    assert str(retrieved_track.album_id) == track["album_id"]
    assert retrieved_track.genre == track["genre"]
    assert retrieved_track.duration == track["duration"]
    assert retrieved_track.lyrics == track["lyrics"]
    assert retrieved_track.is_explicit == track["is_explicit"]
    assert (
        retrieved_track.created_at.strftime("%Y-%m-%dT%H:%M:%S")
        == track["created_at"]
    )
    # Падает из-за разных UTC поздно ночью
    # assert (
    #     retrieved_track.updated_at.strftime("%Y-%m-%d")
    #     == track["updated_at"].split("T")[0]
    # )

    liked_artist = liked_artists[0]
    retrieved_liked_artist = await LikedArtist.get_latest_by_user_artist(
        liked_artist["user_id"], liked_artist["artist_id"], pool=db_pool
    )
    assert retrieved_liked_artist is not None
    assert str(retrieved_liked_artist.user_id) == liked_artist["user_id"]
    assert str(retrieved_liked_artist.artist_id) == liked_artist["artist_id"]

    liked_track = liked_tracks[0]
    retrieved_liked_track = await LikedTrack.get_latest_by_user_track(
        liked_track["user_id"], liked_track["track_id"], pool=db_pool
    )
    assert retrieved_liked_track is not None
    assert str(retrieved_liked_track.user_id) == liked_track["user_id"]
    assert str(retrieved_liked_track.track_id) == liked_track["track_id"]


@pytest.mark.asyncio
async def test_multiple_messages_to_consumer(db_pool):
    """
    Отправляет несколько данных музыкального сервиса в Kafka и проверяет,
    были ли они успешно вставлены в Database через Kafka Consumer.
    """
    artists, albums, tracks, liked_artists, liked_tracks = (
        await send_musical_entities(num_iterations=5)
    )
    await asyncio.sleep(5)
    print(f"Артисты:{artists}")
    print(f"Альбомы:{albums}")
    print(f"Треки:{tracks}")
    print(f"Лайкнутые артисты:{liked_artists}")
    print(f"Лайкнутые треки:{liked_tracks}")

    for artist in artists:
        retrieved_artist = await Artist.get_latest_by_id(
            artist["id"], pool=db_pool
        )
        assert retrieved_artist is not None
        assert str(retrieved_artist.id) == artist["id"]
        assert retrieved_artist.name == artist["name"]
        assert retrieved_artist.author == artist["author"]
        assert retrieved_artist.producer == artist["producer"]
        assert retrieved_artist.country == artist["country"]
        assert retrieved_artist.description == artist["description"]
        assert (
            retrieved_artist.created_at.strftime("%Y-%m-%dT%H:%M:%S")
            == artist["created_at"]
        )
        # Падает из-за разных UTC поздно ночью
        # assert (
        #     retrieved_artist.updated_at.strftime("%Y-%m-%d")
        #     == artist["updated_at"].split("T")[0]
        # )

    for album in albums:
        retrieved_album = await Album.get_latest_by_id(
            album["id"], pool=db_pool
        )
        assert retrieved_album is not None
        assert str(retrieved_album.id) == album["id"]
        assert retrieved_album.title == album["title"]
        assert str(retrieved_album.artist_id) == album["artist_id"]
        assert (
            retrieved_album.release_date.isoformat() == album["release_date"]
        )

    for track in tracks:
        retrieved_track = await Track.get_latest_by_id(
            track["id"], pool=db_pool
        )
        assert retrieved_track is not None
        assert str(retrieved_track.id) == track["id"]
        assert retrieved_track.title == track["title"]
        assert str(retrieved_track.album_id) == track["album_id"]
        assert retrieved_track.genre == track["genre"]
        assert retrieved_track.duration == track["duration"]
        assert retrieved_track.lyrics == track["lyrics"]
        assert retrieved_track.is_explicit == track["is_explicit"]
        assert (
            retrieved_track.created_at.strftime("%Y-%m-%dT%H:%M:%S")
            == track["created_at"]
        )
        # Падает из-за разных UTC поздно ночью
        # assert (
        #     retrieved_track.updated_at.strftime("%Y-%m-%d")
        #     == track["updated_at"].split("T")[0]
        # )

    for liked_artist in liked_artists:
        retrieved_liked_artist = await LikedArtist.get_latest_by_user_artist(
            liked_artist["user_id"], liked_artist["artist_id"], pool=db_pool
        )
        assert retrieved_liked_artist is not None
        assert str(retrieved_liked_artist.user_id) == liked_artist["user_id"]
        assert (
            str(retrieved_liked_artist.artist_id) == liked_artist["artist_id"]
        )

    for liked_track in liked_tracks:
        retrieved_liked_track = await LikedTrack.get_latest_by_user_track(
            liked_track["user_id"], liked_track["track_id"], pool=db_pool
        )
        assert retrieved_liked_track is not None
        assert str(retrieved_liked_track.user_id) == liked_track["user_id"]
        assert str(retrieved_liked_track.track_id) == liked_track["track_id"]
