import time
from tests.fake_producers.fake_auth_producer import send_auth_user

from api.controllers.transform_to_clickhouse_consumer import (
    client
)
from api.models.clickhouse_models import User

from utils.logger_config import configure
import logging
log = logging.getLogger('analysisLogger')
configure()


def test_one_message_to_consumer():
    """
    Отправляет пользовательские данные в Kafka и проверяет,
    были ли они успешно вставлены в ClickHouse через Kafka Consumer.
    """
    user: dict = send_auth_user(num_iterations=1)[0]
    log.debug(f"Sent user: {user}")

    time.sleep(5)

    retrieved_user = User.get_user_by_id(client, user["id"])
    log.info(f"Retrieved user: {retrieved_user}")

    assert retrieved_user is not None
    assert str(retrieved_user.user_id) == user["id"]
    if retrieved_user.gender == "male":
        assert True is user["gender"]
    else:
        assert False is user["gender"]
    assert retrieved_user.age == user["age"]

    format_date = str(retrieved_user.registration_date)
    assert format_date == user["created_at"].split('T')[0]
    assert retrieved_user.citizenship == user["country"]


def test_multiple_messages_to_consumer():
    """
    Отправляет несколько пользовательских данных в Kafka и проверяет,
    были ли они успешно вставлены в ClickHouse через Kafka Consumer.
    """
    users: list[dict] = send_auth_user(num_iterations=5)
    time.sleep(5)

    for user in users:
        retrieved_user = User.get_user_by_id(client, user["id"])
        assert retrieved_user is not None
        assert str(retrieved_user.user_id) == user["id"]
        if retrieved_user.gender == "male":
            assert True is user["gender"]
        else:
            assert False is user["gender"]
        assert retrieved_user.age == user["age"]

        format_date = str(retrieved_user.registration_date)
        assert format_date == user["created_at"].split('T')[0]
        assert retrieved_user.citizenship == user["country"]
