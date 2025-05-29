import asyncio
from app.kafka.consumer import consume
from app.db.connection import connection_check
from app.core.settings import settings

from app.core.logger import get_logger

log = get_logger(__name__)


async def main():
    log.info("Starting Analytics Service...")
    try:
        for attempt in range(10):
            try:
                await connection_check()
                break
            except Exception:
                log.warning(
                    f"""
                    Database not ready yet. Retrying in 5 seconds...
                    (Attempt {attempt + 1}/10)
                """
                )
                await asyncio.sleep(5)
        else:
            log.error(
                f"""
                Failed to connect to Database.
                    host={settings.postgres.host},
                    port={settings.postgres.port},
                    user={settings.postgres.user},
                    password={settings.postgres.password},
                    database={settings.postgres.database},
            """
            )

        for attempt in range(10):
            try:
                await consume()
                break
            except Exception:
                log.warning(
                    f"""
                    Kafka not ready yet. Retrying in 5 seconds...
                    (Attempt {attempt + 1}/10)
                """
                )
                await asyncio.sleep(5)
        else:
            log.error("Failed to connect to Kafka.")
    except Exception as e:
        log.error(f"Analytics Service stopped with error: {e}")
    finally:
        log.info("Analytics Service has been stopped.")


if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        log.info("Analytics Service interrupted by user.")
