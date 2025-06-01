import pytest_asyncio
from asyncpg import create_pool
from app.core.settings import settings


@pytest_asyncio.fixture(scope="function")
async def db_pool():
    pool = await create_pool(
        host=settings.postgres.host,
        port=settings.postgres.port,
        user="user",
        password=settings.postgres.password,
        database=settings.postgres.database,
        min_size=1,
        max_size=5,
    )
    yield pool
    await pool.close()
