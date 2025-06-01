include ../../Makefile

CONFIG_PATH=./configs/test.cfg.yml
COMPOSE_PATH=./configs/docker-compose.yml
MIGRATIONS_DIR=./migrations
DB_URL=postgresql://user:password@localhost:5432/music
TEST_COMPOSE_PATH=./configs/test.docker-compose.yml
TEST_DB_URL=postgres://test:test@localhost:5432/test-musical?sslmode=disable

.PHONY: build
build: gen-pb
	@go build -o bin/musical cmd/musical/main.go

.PHONY: run
run: build
	bin/musical --path $(CONFIG_PATH)

# Docker
.PHONY: compose-up
compose-up:
	@docker compose -f $(COMPOSE_PATH) up -d
	@docker compose -f $(COMPOSE_PATH) exec -T postgres bash -c \
		'until pg_isready -U user -d music; do sleep 1; done'

.PHONY: compose-down
compose-down:
	@docker compose -f $(COMPOSE_PATH) down 

.PHONY: test-compose-up
test-compose-up:
	@docker compose -f $(TEST_COMPOSE_PATH) up -d
	@docker compose -f $(TEST_COMPOSE_PATH) exec -T postgres bash -c \
		'until pg_isready -U user -d music; do sleep 1; done'

.PHONY: test-compose-down
test-compose-down:
	@docker compose -f $(TEST_COMPOSE_PATH) down 

# Migrations
.PHONY: migrate-up
migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

.PHONY: migrate-down
migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

.PHONY: test-migrate-up
test-migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(TEST_DB_URL)" up

.PHONY: test-migrate-down
test-migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(TEST_DB_URL)" down

# Protobuf
.PHONY: gen-pb
gen-pb:
	@protoc --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative \
		    --go_out=pkg/musicalpb --go-grpc_out=pkg/musicalpb \
		    -I api api/*

# Mocks
.PHONY: gen-mock
gen-mock:
	@find . -name '*_mock.go' -delete
	mockgen -destination=tests/mocks/album_repo_mock.go -package=mocks github.com/TeoPlow/online-music-service/src/musical/internal/domain AlbumRepo
	mockgen -destination=tests/mocks/track_repo_mock.go -package=mocks github.com/TeoPlow/online-music-service/src/musical/internal/domain TrackRepo
	mockgen -destination=tests/mocks/artist_service_mock.go -package=mocks github.com/TeoPlow/online-music-service/src/musical/internal/domain ArtistClient
	mockgen -destination=tests/mocks/album_service_mock.go -package=mocks github.com/TeoPlow/online-music-service/src/musical/internal/domain AlbumClient
	mockgen -destination=tests/mocks/streaming_service_mock.go -package=mocks github.com/TeoPlow/online-music-service/src/musical/internal/domain StreamingClient
	mockgen -destination=tests/mocks/tx_manager_mock.go -package=mocks github.com/TeoPlow/online-music-service/src/musical/internal/domain TxManager

# for tests
.PHONY: test-deps-up
test-deps-up: test-compose-up test-migrate-up

.PHONY: test-deps-down
test-deps-down: test-compose-down