from asyncpg import create_pool
from app.core.settings import settings

from app.core.logger import get_logger

log = get_logger(__name__)

_db_pool = None


async def get_db_client():
    """
    Возвращает подключение к клиенту PostgreSQL через asyncpg.
    """
    global _db_pool

    if _db_pool is None:
        _db_pool = await create_pool(
            host=settings.postgres.host,
            port=settings.postgres.port,
            user="postgres",
            password=settings.postgres.password,
            database=settings.postgres.database,
            min_size=1,
            max_size=10,
        )
    return _db_pool


async def connection_check():
    log.info("Checking connection...")

    client = await get_db_client()
    async with client.acquire() as conn:
        await conn.execute(
            """
        SELECT 1;
        """
        )

        log.info("Connection is successful.")
