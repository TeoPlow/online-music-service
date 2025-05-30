FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apk add --no-cache protobuf && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    go build -o auth_service ./src/auth_service/cmd/auth_service/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/auth_service .

COPY src/auth_service/config ./config
COPY src/auth_service/.env .

EXPOSE 50051

CMD ["./auth_service", "--path", "./config/config.yaml"]
