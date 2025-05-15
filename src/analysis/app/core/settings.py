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
    clickhouse_host = "localhost"
else:
    kafka_bootstrap = "kafka:9092"
    clickhouse_host = "clickhouse"


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
        "music-liked-tracks"
    ]


class ClickHouseConfig(BaseSettings):
    host: str = Field(default_factory=lambda: clickhouse_host)
    port: int = 9000
    user: str = "default"
    password: str = 'mysecurepassword'
    database: str = "music_streaming"


class AppConfig(BaseSettings):
    base_dir: Path = Path(__file__).resolve().parent.parent.parent


class Settings(BaseSettings):
    log: LogConfig = LogConfig()
    paths: AppConfig = AppConfig()
    kafka: KafkaConfig = KafkaConfig()
    clickhouse: ClickHouseConfig = ClickHouseConfig()

    model_config = ConfigDict(
        env_nested_delimiter="__",
        case_sensitive=False
    )


yaml_config = load_yaml_config()
settings = Settings(**yaml_config)
