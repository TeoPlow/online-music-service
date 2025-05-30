include ../../Makefile

MIGRATIONS_DIR=./app/db/migrations/sql
CONFIG_PATH=./config.yml
COMPOSE_PATH=./docker-compose.yml
DB_URL=postgres://user:password@localhost:5432/analysis_db
TESTS_PATH=./tests

.PHONY: setup-venv start-services stop-services remove-volumes postgres-shell

# Запуск сервисов в фоне
run-detached:
	@echo "Starting services..."
	docker-compose -f ${COMPOSE_PATH} up --build -d

# Запуск сервисов в интерактивном режиме (в терминале)
run:
	@echo "Starting services..."
	docker-compose -f ${COMPOSE_PATH} up --build

# Остановка сервисов
stop-services:
	@echo "Stopping services..."
	docker-compose -f ${COMPOSE_PATH} down

# Запуск миграций
.PHONY: migrate-up
migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

# Удаление всех томов и остановка сервисов
remove-volumes:
	@echo "Removing volumes..."
	docker-compose -f ${COMPOSE_PATH} down --volumes

# Открытие оболочки Postgres
postgres-shell:
	@echo "Opening PostgreSQL shell..."
	docker exec -it analysis-postgres-1 psql -U postgres

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