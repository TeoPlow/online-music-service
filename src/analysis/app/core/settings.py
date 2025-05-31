import os
from pydantic_settings import BaseSettings
from pydantic import ConfigDict, Field
from typing import List
from pathlib import Path
from app.core.config import load_yaml_config

# Настройка для тестирования
env = os.getenv("APP_ENV", "default")
if env == "test":
    kafka_bootstrap = "localhost:9094"
    db_host = "localhost"
else:
    kafka_bootstrap = "kafka:9092"
    db_host = "postgres"


class LogConfig(BaseSettings):
    level: str = "INFO"
    save_to_file: bool = False


class KafkaConfig(BaseSettings):
    bootstrap_servers: str = Field(default_factory=lambda: kafka_bootstrap)
    topics: List[str] = [
        "auth-users",
        "music-artists",
        "music-albums",
        "music-tracks",
        "music-liked-artists",
        "music-liked-tracks",
    ]


class PostgresConfig(BaseSettings):
    host: str = Field(default_factory=lambda: db_host)
    port: int = 5432
    user: str = "user"
    password: str = "password"
    database: str = "analysis_db"


class AppConfig(BaseSettings):
    base_dir: Path = Path(__file__).resolve().parent.parent.parent


class Settings(BaseSettings):
    log: LogConfig = LogConfig()
    paths: AppConfig = AppConfig()
    kafka: KafkaConfig = KafkaConfig()
    postgres: PostgresConfig = PostgresConfig()

    model_config = ConfigDict(env_nested_delimiter="__", case_sensitive=False)


yaml_config = load_yaml_config()
settings = Settings(**yaml_config)
