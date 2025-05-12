include ../../Makefile

CONFIG_PATH=./config.yml
COMPOSE_PATH=./configs/docker-compose.yml
TESTS_PATH=./tests
MIGRATIONS_DIR=./migrations
PYTHONPATH=.
API_PATH=./api

.PHONY: setup-venv set-pythonpath start-services stop-services remove-volumes \
	run-fake-auth-producer run-fake-musical-producer run-consumer run-tests \
	clickhouse-shell

# Полный запуск сервиса с нуля
full-run: setup-venv start-services 
	@echo "Wait for 7 seconds to start all services..."
	sleep 7
	
	@echo "Running consumer..."
	PYTHONPATH=. .venv/bin/python ${API_PATH}/controllers/transform_to_clickhouse_consumer.py


# Тут перечислены все цели, которые можно запустить из Makefile
setup-venv:
	@echo "Setting up virtual environment..."
	python -m venv .venv

	@echo "Installing dependencies..."
	.venv/bin/pip install -r ../../requirements.txt
	.venv/bin/pip freeze 

set-pythonpath:
	@echo "Set pythonpath"
	export PYTHONPATH=${PYTHONPATH}

start-services:
	@echo "Starting services..."
	docker-compose -f ${COMPOSE_PATH} up -d

stop-services:
	@echo "Stopping services..."
	docker-compose -f ${COMPOSE_PATH} down

remove-volumes:
	@echo "Removing volumes..."
	docker-compose -f ${COMPOSE_PATH} down --volumes

run-fake-auth-producer:
	@echo "Running fake auth producer..."
	PYTHONPATH=. .venv/bin/python ${TESTS_PATH}/fake_producers/fake_auth_producer.py

run-fake-musical-producer:
	@echo "Running fake musical producer..."
	PYTHONPATH=. .venv/bin/python ${TESTS_PATH}/fake_producers/fake_musical_producer.py

run-consumer:
	@echo "Running consumer..."
	PYTHONPATH=. .venv/bin/python ${API_PATH}/controllers/transform_to_clickhouse_consumer.py

# Запускать после `make full-run` в отельном терминале
run-tests:
	@echo "Running tests..."
	PYTHONPATH=. .venv/bin/pytest ${TESTS_PATH}/

clickhouse-shell:
	@echo "Opening ClickHouse shell..."
	docker exec -it clickhouse-server clickhouse-client