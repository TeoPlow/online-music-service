include ../../Makefile

CONFIG_PATH=./config.yml
COMPOSE_PATH=./docker-compose.yml
TESTS_PATH=./tests
SUPERSET_COMPOSE_PATH=./superset/docker-compose-image-tag.yml
MIGRATIONS_DIR=./app/db/migrations/sql
DB_URL=postgres://user:password@localhost:5432/analysis_db

.PHONY: setup-venv start-services stop-services remove-volumes postgres-shell

# Создание общей docker-сети
create-network:
	@echo "Creating Docker network 'music_net' if it doesn't exist..."
	docker network create music_net || true

# Запуск сервисов в фоне
start-services:
	@echo "Starting services..."
	docker compose -f ${COMPOSE_PATH} up --build -d
	@echo "Starting Superset services..."
	docker compose -f ${SUPERSET_COMPOSE_PATH} up --build -d

start-services+:
	@echo "Starting services..."
	docker compose -f ${COMPOSE_PATH} up --build

# Остановка сервисов
stop-services:
	@echo "Stopping services..."
	docker compose -f ${COMPOSE_PATH} down
	@echo "Stopping superset..."
	docker compose -f ${SUPERSET_COMPOSE_PATH} down

# Удаление всех томов и остановка сервисов
remove-volumes:
	@echo "Removing volumes..."
	docker compose -f ${COMPOSE_PATH} down --volumes

# Открытие оболочки Postgres
postgres-shell:
	@echo "Opening PostgreSQL shell..."
	docker exec -it analysis-postgres-1 psql -U postgres

# Вывод всех логов
full-logs:
	@echo "Viewing service logs..."
	docker compose -f ${COMPOSE_PATH} logs --tail 100 -f

# Вывод логов сервиса
service-logs:
	@echo "Viewing service logs..."
	docker logs analysis

run-tests:
	@echo "Starting tests..."
	pytest ${TESTS_PATH}

# Миграции
.PHONY: migrate-up
migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

# Для системы CI Евгеновской
test-deps-up: start-services
test-deps-down: remove-volumes
