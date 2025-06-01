FROM golang:1.24.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apk add --no-cache protobuf && \
    go build -o auth ./src/auth/cmd/auth/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/auth .

EXPOSE 50051

CMD ["./auth", "--path", "/etc/auth/configs/config.yaml"]