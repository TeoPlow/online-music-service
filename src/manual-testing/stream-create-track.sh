#!/bin/bash

# Конфигурация
HOST="localhost:8080"
PROTOSET="descriptor.pb"
SERVICE="musical.MusicalService"
AUTH_TOKEN=$(cat token.txt)
FILE_PATH="../musical/tests/testdata/meow-meow.mp3"

REQUEST_FILE=$(mktemp)

# 1. Сначала добавляем метаданные
cat <<EOF > $REQUEST_FILE
{
  "metadata": {
    "title": "Meow Meow Song",
    "album_id": "bf13d0c8-7ac7-4806-a4c2-fb547e8e29de",
    "genre": "pop",
    "duration": 237,
    "lyrics": "Meow meow meow meow...",
    "is_explicit": false
  }
}
EOF

# 2. Затем добавляем бинарные данные файла
echo -n '{"chunk": "' >> $REQUEST_FILE
base64 -w0 $FILE_PATH >> $REQUEST_FILE
echo '"}' >> $REQUEST_FILE

grpcurl -plaintext \
  -protoset $PROTOSET \
  -rpc-header "authorization: Bearer $AUTH_TOKEN" \
  -d @ $HOST $SERVICE/CreateTrack < $REQUEST_FILE

rm $REQUEST_FILE