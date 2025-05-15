from asynch import Connection
from app.core.settings import settings


async def get_clickhouse_client() -> Connection:
    """
    Возвращает подключение к клиенту ClickHouse через asynch.
    """
    conn = Connection(
        host=settings.clickhouse.host,
        port=settings.clickhouse.port,
        user="default",
        # user=settings.clickhouse.user,
        password=settings.clickhouse.password,
        database="default",
    )
    await conn.connect()
    return conn
