import time
import random
import uuid
from datetime import datetime
from kafka import KafkaProducer
import json
from faker import Faker

from configs.config import (
    kafka_bootstrap_servers,
)

from utils.logger_config import configure
import logging
log = logging.getLogger('analysisLogger')
configure()

fake = Faker()

log.debug("[auth-users] Запущен фейковый продюсер.")

producer = KafkaProducer(
    bootstrap_servers=kafka_bootstrap_servers,
    value_serializer=lambda v: json.dumps(v).encode('utf-8')
)

roles = ['user', 'admin', 'artist']

with open('../static/countries.json', 'r') as f:
    countries = json.load(f)


def send_auth_user(num_iterations=1):
    users = []
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
            "created_at": datetime.now().isoformat()
        }
        producer.send('auth-users', user)
        log.debug(f"[auth-users] Kafka Sent: {user}")
        users.append(user)
    return users


if __name__ == '__main__':
    while True:
        _ = send_auth_user()
        time.sleep(random.randint(5, 10))
