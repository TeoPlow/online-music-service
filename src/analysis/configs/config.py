import yaml  # type: ignore
import os


def get_base_dir():
    path = os.path.abspath(os.path.join(os.path.dirname(__file__), '..'))
    return path


config = {}

try:
    config_path = get_base_dir() + "/config.yml"
    with open(config_path) as f:
        config = yaml.safe_load(f)
except FileNotFoundError:
    print("Config file not found on path:", config_path)

# Пути
BASE_DIR = config.get("paths", {}).get("base_dir", get_base_dir())

# Логирование
LOG_LEVEL = config.get("log", {}).get("level", "INFO").upper()
SAVE_LOG_TO_FILE = config.get("log", {}).get("save_to_file", False)

# Kafka
kafka_bootstrap_servers = config.get("kafka", {}).get(
    "bootstrap_servers", "localhost:9092"
    )
kafka_topics = config.get("kafka", {}).get("topics", [])

# ClickHouse
clickhouse_host = config.get("clickhouse", {}).get("host", "localhost")
clickhouse_port = config.get("clickhouse", {}).get("port", 9000)
clickhouse_user = config.get("clickhouse", {}).get("user", "default")
clickhouse_password = config.get("clickhouse", {}).get(
    "password",
    ""
)
clickhouse_database = config.get("clickhouse", {}).get(
    "database",
    "music_streaming"
)


if __name__ == "__main__":
    print(f"BASE_DIR: {BASE_DIR}")
