# Запуск Analysis
1. Создать `.venv`
2. Прописать эту шнягу в терминал (если чёто не работает): 
`export PYTHONPATH=path/to/online-music-service`
3. Запустить ClickHouse, Kafka и Zookeeper
4. Запустить скрипты лже-продюсеров - `fake_auth_producer.py` и `fake_music_producer.py`
5. Запустить скрипт консьюмера - `transfromation_to_clickhouse.py`


# Запуск ClickHouse, Kafka и Zookeeper для тестирования
`docker-compose -f src/analysis/tests/docker-compose.test.yml up -d`


### Тут можно делать команды в ClickHouse (например, чекнуть БД)
`docker exec -it clickhouse-server clickhouse-client`


# Запуск тестов
`pytest src/analysis/tests/`
