include ../../Makefile

CONFIG_PATH=./configs/cfg.yml
COMPOSE_PATH=./configs/docker-compose.yml
MIGRATIONS_DIR=./migrations
DB_URL=postgresql://user:password@localhost:5432/music

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

.PHONY: compose-down
compose-down:
	@docker compose -f $(COMPOSE_PATH) down 

# Migrations
.PHONY: migrate-up
migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

.PHONY: migrate-down
migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

# Protobuf
.PHONY: gen-pb
gen-pb:
	@protoc --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative \
		    --go_out=pkg/musicalpb --go-grpc_out=pkg/musicalpb \
		    -I api api/musical.proto
