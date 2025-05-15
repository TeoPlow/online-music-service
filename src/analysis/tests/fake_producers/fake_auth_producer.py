import random
import uuid
from datetime import datetime
from aiokafka import AIOKafkaProducer
import json
from faker import Faker
from app.core.settings import settings

from app.kafka.topics import (
    AUTH_USERS,
)

from app.core.logger import get_logger
log = get_logger(__name__)

fake = Faker()

roles = ['user', 'admin', 'artist']


with open(f'{settings.paths.base_dir}/../static/countries.json', 'r') as f:
    countries = json.load(f)


async def send_auth_user(num_iterations=1):
    """
    Это функция с запуском фейкового асинхронного Kafka Producer,
    который будет спамить заданное кол-во раз автосгенерированного
    пользователя.
    """
    users = []
    producer = AIOKafkaProducer(
        bootstrap_servers=settings.kafka.bootstrap_servers,
        value_serializer=lambda v: json.dumps(v).encode('utf-8')
    )
    await producer.start()
    try:
        for _ in range(num_iterations):
            user = {
                "id": str(uuid.uuid4()),
                "username": fake.user_name(),
                "email": fake.email(),
                "gender": random.choice([True, False]),
                "country": random.choice(countries),
                "age": random.randint(14, 80),
                "role": random.choice(roles),
                "passHash": fake.sha256(),
                "created_at": datetime.now().isoformat().split('.')[0]
            }
            await producer.send(AUTH_USERS, user)
            log.debug(f"[{AUTH_USERS}] Kafka Sent: {user}")
            users.append(user)
    finally:
        await producer.stop()
    return users


if __name__ == '__main__':
    import asyncio
    asyncio.run(send_auth_user())
