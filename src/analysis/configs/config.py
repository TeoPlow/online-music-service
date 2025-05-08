import yaml  # type: ignore

with open("src/analysis/configs/config.yml") as f:
    config = yaml.safe_load(f)


# Логирование
log_level = config.get("log", {}).get("level", "INFO").upper()
save_to_file = config.get("log", {}).get("save_to_file", False)

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
