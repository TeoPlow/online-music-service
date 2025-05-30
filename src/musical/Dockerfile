FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apk add --no-cache protobuf git && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    go build -o musical ./src/musical/cmd/musical/main.go


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/musical .

COPY src/musical/configs ./configs

EXPOSE 50052

CMD ["./musical", "--path", "./configs/cfg.yml"]
