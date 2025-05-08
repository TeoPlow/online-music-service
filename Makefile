PYENV=python -m 


# Go tasks
.PHONY: go-lint
go-lint:
	@echo "Running Go linter..."
	@if [ -z "$$(find . -type f -name '*.go')" ]; then \
		echo "No Go files found."; \
	else \
		golangci-lint run; \
	fi

.PHONY: go-deps
go-deps:
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/golang/mock/mockgen@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go mod tidy

.PHONY: go-fmt
go-fmt:
	gofumpt -w .

.PHONY: go-test
go-test:
	@echo "Running Go tests..."
	@if [ -z "$$(find . -type f -name '*.go')" ]; then \
		echo "No Go tests found."; \
	else \
		go test -v ./...; \
	fi


# Python tasks
.PHONY: py-lint
py-lint:
	$(PYENV) flake8 . --exclude .venv

.PHONY: py-deps
py-deps:
	$(PYENV) pip install --upgrade pip
	$(PYENV) pip install flake8 pytest black
	$(PYENV) pip install -r requirements.txt

.PHONY: py-fmt
py-fmt:
	$(PYENV) black . --exclude .venv

.PHONY: py-test
py-test:
	@echo "Running Python tests..."
	@echo "Starting docker-compose for test dependencies..."
	@docker-compose -f src/analysis/tests/docker-compose.test.yml up -d
	@echo "Waiting for ClickHouse to be ready..."
	@for i in 1 2 3 4 5; do \
		nc -z localhost 9000 && echo "ClickHouse is up!" && break || \
		(echo "Waiting..."; sleep 2); \
	done
	@if [ -z "$$(find . -type f -name 'test_*.py' -o -name '*_test.py')" ]; then \
		echo "No Python tests found."; \
	else \
		pytest; \
	fi

# Common tasks
.PHONY: lint
lint: go-lint py-lint

.PHONY: deps
deps: go-deps py-deps

.PHONY: fmt
fmt: go-fmt py-fmt

.PHONY: test
test: go-test py-test
