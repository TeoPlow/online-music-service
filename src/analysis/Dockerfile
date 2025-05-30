FROM python:3.12-slim

ENV PYTHONPATH=/app
WORKDIR /app

COPY ../../requirements.txt .

RUN apt-get update && apt-get install -y gcc pkg-config

RUN pip install --no-cache-dir -r requirements.txt

COPY src/analysis/. .

CMD ["python", "app/main.py"]
