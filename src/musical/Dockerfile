FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apk add --no-cache protobuf git && \
    go build -o musical ./src/musical/cmd/musical/main.go


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/musical .

EXPOSE 50052

CMD ["./musical", "--path", "/etc/musical/configs/config.yaml"]