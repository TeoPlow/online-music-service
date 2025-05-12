from utils.logger_config import configure
import logging
log = logging.getLogger('analysisLogger')
configure()


class User:
    def __init__(self,
                 user_id,
                 gender,
                 age,
                 registration_date,
                 citizenship,
                 liked_songs=None,
                 liked_artists=None):
        self.user_id = user_id
        self.gender = gender
        self.age = age
        self.registration_date = registration_date
        self.citizenship = citizenship
        self.liked_songs = liked_songs or []
        self.liked_artists = liked_artists or []

    def save(self, client):
        client.execute(
            """
            INSERT INTO music_streaming.users (
                user_id,
                gender,
                age,
                registration_date,
                citizenship,
                liked_songs,
                disliked_songs,
                liked_artists,
                disliked_artists
            ) VALUES
            """,
            [
                (
                    self.user_id,
                    self.gender,
                    self.age,
                    self.registration_date,
                    self.citizenship,
                    self.liked_songs,
                    [],  # Всё сложно с disliked_songs...
                    self.liked_artists,
                    [],  # Всё сложно с disliked_artists...
                )
            ]
        )
        log.info("User inserted successfully: %s", self.user_id)

    @classmethod
    def get_user_by_id(cls, client, user_id):
        """
        Получает пользователя из ClickHouse по его ID.
        """
        result = client.execute(
            """
            SELECT * FROM music_streaming.users
            WHERE user_id = toUUID(%(user_id)s)
            """,
            {'user_id': user_id}
        )
        if result:
            log.debug(f"User found in ClickHouse: {result[0]}")
            user_data = {
                'user_id': result[0][0],
                'gender': result[0][1],
                'age': result[0][2],
                'registration_date': result[0][3],
                'citizenship': result[0][4],
                'liked_songs': result[0][5],
                'liked_artists': result[0][6],
            }
            return cls(**user_data)
        return None


class Artist:
    def __init__(self, artist_id, artist_name, citizenship):
        self.artist_id = artist_id
        self.artist_name = artist_name
        self.citizenship = citizenship

    def save(self, client):
        client.execute(
            """
            INSERT INTO music_streaming.artists (
                artist_id,
                artist_name,
                citizenship
            ) VALUES
            """,
            (
                self.artist_id,
                self.artist_name,
                self.citizenship
            )
        )
        log.info("Artist inserted successfully: %s", self.artist_id)

    @classmethod
    def get_artist_by_id(cls, client, artist_id):
        """
        Получает артиста из ClickHouse по его ID.
        """
        result = client.execute(
            """
            SELECT
                artist_id,
                artist_name,
                citizenship
            FROM music_streaming.artists
            WHERE artist_id = toUUID(%(artist_id)s)
            """,
            {'artist_id': artist_id}
        )
        if result:
            log.debug("Artist found in ClickHouse: %s", result[0])
            artist_data = {
                'artist_id': result[0][0],
                'artist_name': result[0][1],
                'citizenship': result[0][2],
            }
            return cls(**artist_data)
        return None


class Song:
    def __init__(self, song_id, song_name, genre, artist_id):
        self.song_id = song_id
        self.song_name = song_name
        self.genre = genre
        self.artist_id = artist_id

    def save(self, client):
        client.execute(
            """
            INSERT INTO music_streaming.songs (
                song_id,
                song_name,
                genre,
                artist_id
            ) VALUES
            """,
            (
                self.song_id,
                self.song_name,
                self.genre,
                self.artist_id
            )
        )
        log.info("Song inserted successfully: %s", self.song_id)

    @classmethod
    def get_song_by_id(cls, client, song_id):
        """
        Получает песню из ClickHouse по её ID.
        """
        result = client.execute(
            """
            SELECT
                song_id,
                song_name,
                genre,
                artist_id
            FROM music_streaming.songs
            WHERE song_id = toUUID(%(song_id)s)
            """,
            {'song_id': song_id}
        )
        if result:
            log.debug(f"Song found in ClickHouse: {result[0]}")
            song_data = {
                'song_id': result[0][0],
                'song_name': result[0][1],
                'genre': result[0][2],
                'artist_id': result[0][3],
            }
            return cls(**song_data)
        return None


class LikedArtist:
    def __init__(self, user_id, artist_id):
        self.user_id = user_id
        self.artist_id = artist_id

    def save(self, client):
        client.execute(
            """
            INSERT INTO music_streaming.liked_artists (
                user_id,
                artist_id
            ) VALUES
            """,
            (
                self.user_id,
                self.artist_id
            )
        )
        log.info("Liked artist added successfully: %s", self.artist_id)

    @classmethod
    def get_liked_artist_by_id(cls, client, liked_artist_id):
        """
        Получает информацию о любимом артисте из ClickHouse по его ID.
        """
        result = client.execute(
            """
            SELECT
                user_id,
                artist_id
            FROM music_streaming.liked_artists
            WHERE artist_id = toUUID(%(liked_artist_id)s)
            """,
            {'liked_artist_id': liked_artist_id}
        )
        if result:
            log.debug(f"Liked artist found in ClickHouse: {result[0]}")
            liked_artist_data = {
                'user_id': result[0][0],
                'artist_id': result[0][1],
            }
            return cls(**liked_artist_data)
        return None


class LikedTrack:
    def __init__(self, user_id, track_id):
        self.user_id = user_id
        self.track_id = track_id

    def save(self, client):
        client.execute(
            """
            INSERT INTO music_streaming.liked_tracks (
                user_id,
                track_id
            ) VALUES
            """,
            (
                self.user_id,
                self.track_id
            )
        )
        log.info("Liked track added successfully: %s", self.track_id)

    @classmethod
    def get_liked_track_by_id(cls, client, liked_track_id):
        """
        Получает информацию о любимой треке из ClickHouse по его ID.
        """
        result = client.execute(
            """
            SELECT
                user_id,
                track_id
            FROM music_streaming.liked_tracks
            WHERE track_id = toUUID(%(liked_track_id)s)
            """,
            {'liked_track_id': liked_track_id}
        )
        if result:
            log.debug(f"Liked track found in ClickHouse: {result[0]}")
            liked_track_data = {
                'user_id': result[0][0],
                'track_id': result[0][1],
            }
            return cls(**liked_track_data)
        return None
