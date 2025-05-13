SERVICES := musical analysis # Сюда пишем название папки сервиса через пробел
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

# Common tasks
.PHONY: lint
lint: go-lint py-lint

.PHONY: deps
deps: go-deps py-deps

.PHONY: fmt
fmt: go-fmt py-fmt

.PHONY: test
test:
	@if [ -n "$(SERVICE)" ]; then \
		echo "Running tests only for $(SERVICE)..."; \
		echo "=> Starting dependencies for $(SERVICE)..."; \
		(cd ./src/$(SERVICE) && $(MAKE) -s test-deps-up) || echo "No test-deps-up for $(SERVICE)"; \
		echo "=> Running Go tests..."; \
		if find ./src/$(SERVICE) -type f -name '*.go' | grep -q .; then \
			go test -v ./src/$(SERVICE)/... || exit 1; \
		else \
			echo "No Go files found in ./src/$(SERVICE)."; \
		fi; \
		echo "=> Running Python tests..."; \
		if find ./src/$(SERVICE) -type f -name '*_test.py' | grep -q .; then \
			pytest ./src/$(SERVICE) || exit 1; \
		else \
			echo "No Python tests found in ./src/$(SERVICE)."; \
		fi; \
		echo "=> Stopping dependencies..."; \
		(cd ./src/$(SERVICE) && $(MAKE) -s test-deps-down) || echo "No test-deps-down for $(SERVICE)"; \
	else \
		echo "Running tests for all services..."; \
		for service in $(SERVICES); do ( \
			echo "=> Starting dependencies for $$service..."; \
			(cd ./src/$$service && $(MAKE) -s test-deps-up) || echo "No test-deps-up for $$service"; \
			echo "=> Running Go tests for $$service..."; \
			if find ./src/$$service -type f -name '*.go' | grep -q .; then \
				go test -v ./src/$$service/... || exit 1; \
			else \
				echo "No Go files found in ./src/$$service."; \
			fi; \
			echo "=> Running Python tests for $$service..."; \
			if find ./src/$$service -type f -name '*_test.py' | grep -q .; then \
				pytest ./src/$$service || exit 1; \
			else \
				echo "No Python tests found in ./src/$$service."; \
			fi; \
			echo "=> Stopping dependencies for $$service..."; \
			(cd ./src/$$service && $(MAKE) -s test-deps-down) || echo "No test-deps-down for $$service"; \
		); done; \
	fi

.PHONY: help
help:
	@echo "Available make targets:"
	@echo
	@echo "  help          Show this help message"
	@echo
	@echo "Go tasks:"
	@echo "  go-lint       Run Go linter (golangci-lint)"
	@echo "  go-deps       Install Go dependencies (linters, tools, etc.)"
	@echo "  go-fmt        Format Go code using gofumpt"
	@echo
	@echo "Python tasks:"
	@echo "  py-lint       Run Python linter (flake8)"
	@echo "  py-deps       Install Python development dependencies"
	@echo "  py-fmt        Format Python code using black"
	@echo
	@echo "Common tasks:"
	@echo "  lint          Run both Go and Python linters"
	@echo "  deps          Install both Go and Python dependencies"
	@echo "  fmt           Format both Go and Python code"
	@echo
	@echo "Test tasks:"
	@echo "  test          Run tests for all services or a specific one"
	@echo "                Usage: make test SERVICE=<service-name>"
	@echo "                It is required that the commands 'make test-deps-up'"
	@echo "                and 'make test-deps-down' are executed inside the service folder"
