import asyncio
import pytest
from app.db.models import User
from tests.fake_producers.fake_auth_producer import send_auth_user
from app.core.settings import settings


@pytest.mark.asyncio
async def test_one_message_to_consumer():
    """
    Отправляет пользовательские данные в Kafka и проверяет,
    были ли они успешно вставлены в ClickHouse через Kafka Consumer.
    """
    ch = settings.clickhouse
    print(
        f"[CONFIG] ClickHouse settings:\n"
        f"  host={ch.host}\n"
        f"  port={ch.port}\n"
        f"  user={ch.user}\n"
        f"  password={ch.password[:4]+'***' if ch.password else '(empty)'}\n"
    )

    user: dict = (await send_auth_user(num_iterations=1))[0]
    print(f"Sent user: {user}")

    await asyncio.sleep(5)

    retrieved_user: User = await User.get_latest_by_id(user["id"])
    print(f"Retrieved user: {retrieved_user}")

    assert retrieved_user is not None
    assert str(retrieved_user.id) == user["id"]
    if retrieved_user.gender == "male":
        assert True is user["gender"]
    else:
        assert False is user["gender"]
    assert retrieved_user.age == user["age"]

    # Миллисекунды ClickHouse не хранит (вроде)
    date = retrieved_user.created_at.strftime('%Y-%m-%dT%H:%M:%S')
    assert date == user["created_at"]
    assert retrieved_user.country == user["country"]


@pytest.mark.asyncio
async def test_multiple_messages_to_consumer():
    """
    Отправляет несколько пользовательских данных в Kafka и проверяет,
    были ли они успешно вставлены в ClickHouse через Kafka Consumer.
    """
    users: list[dict] = await send_auth_user(num_iterations=5)
    await asyncio.sleep(5)

    for user in users:
        retrieved_user: User = await User.get_latest_by_id(user["id"])
        print(f"Retrieved user: {retrieved_user}")
        assert retrieved_user is not None
        assert str(retrieved_user.id) == user["id"]
        if retrieved_user.gender == "male":
            assert True is user["gender"]
        else:
            assert False is user["gender"]
        assert retrieved_user.age == user["age"]

        # Миллисекунды ClickHouse не хранит (вроде)
        date = retrieved_user.created_at.strftime('%Y-%m-%dT%H:%M:%S')
        assert date == user["created_at"]
        assert retrieved_user.country == user["country"]


# Написать тесты для фейкового продюсера из музыкального сервиса.
