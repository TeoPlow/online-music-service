FROM python:3.12-slim

ENV PYTHONPATH=/app
WORKDIR /app

COPY ../../requirements.txt .

RUN apt-get update && apt-get install -y gcc pkg-config curl

RUN pip install --no-cache-dir -r requirements.txt

COPY src/analysis/. .

ENV GOPATH=/go PATH=$PATH:$GOPATH/bin:/usr/local/go/bin

RUN curl -L https://go.dev/dl/go1.21.0.linux-amd64.tar.gz | tar xz -C /usr/local && \
    go install github.com/pressly/goose/v3/cmd/goose@latest

CMD ["sh", "-c", "goose -dir /app/db/migrations/sql postgres 'postgresql://user:password@postgres:5432/analysis_db' up && python app/main.py"]
