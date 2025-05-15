include ../../Makefile

CONFIG_PATH=./config.yml
COMPOSE_PATH=./docker-compose.yml
TESTS_PATH=./tests

.PHONY: setup-venv start-services stop-services remove-volumes clickhouse-shell

# Запуск сервисов в фоне
start-services:
	@echo "Starting services..."
	docker-compose -f ${COMPOSE_PATH} up --build -d

# Остановка сервисов
stop-services:
	@echo "Stopping services..."
	docker-compose -f ${COMPOSE_PATH} down

# Удаление всех томов и остановка сервисов
remove-volumes:
	@echo "Removing volumes..."
	docker-compose -f ${COMPOSE_PATH} down --volumes

# Открытие оболочки ClickHouse
clickhouse-shell:
	@echo "Opening ClickHouse shell..."
	docker exec -it clickhouse-server clickhouse-client

# Вывод всех логов
full-logs:
	@echo "Viewing service logs..."
	docker-compose -f ${COMPOSE_PATH} logs --tail 100 -f

# Вывод логов сервиса
service-logs:
	@echo "Viewing service logs..."
	docker logs analysis

run-tests:
	@echo "Starting tests..."
	pytest ${TESTS_PATH}

# Для системы CI Евгеновской
test-deps-up: start-services
test-deps-down: remove-volumes