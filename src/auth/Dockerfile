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

RUN mkdir -p /app/static
COPY --from=builder /app/src/static /app/static
ENV STATIC_FILES_PATH=/app/static
RUN mkdir -p /etc/auth/configs

COPY --from=builder /app/src/auth/config/config.yaml /etc/auth/configs/

RUN ls -la /app/static/
RUN ls -la /etc/auth/configs/

EXPOSE 50051

CMD ["./auth", "--path", "/etc/auth/configs/config.yaml"]