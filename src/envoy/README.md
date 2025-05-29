Пока запуск такой:
1. Запустить оба сервиса (auth_service и musical_service) на локальной машине.
    Делается это так:
    - Musical service: 
        1. Добавить `export PATH="$PATH:$(go env GOPATH)/bin"` в конец файла ~/.zshrc
        2. В терминале прописать `source ~/.zshrc`
        3. Поменять gRPC порт на `grpc_port: :50052`
        4. Запустить БД командой `make compose-up` 
        5. Запустить миграции командой `make migrate-up` 
        6. Запустить сервис командой `make run`
    
    - Auth service: 
        1. В терминале прописать `source ~/.zshrc`
        2. Добавть .env файл в корне auth_service с писаниной: `CONFIG_PATH = ./config/config.yaml`
        3. Поменять порты у БД и у Redis на свободные (например, порт для Redis: 6378 и для БД: 5433)
        4. Запустить БД командой `make compose-up` 
        5. Запустить миграции командой `make migrate-up` 
        6. Запустить сервис командой `make run`

2. Запускаем envoy из папки src/envoy:
    ```
    docker run --rm -it -p 8080:8080 -p 9901:9901 \
        -v $(pwd)/envoy.yaml:/etc/envoy/envoy.yaml \
        envoyproxy/envoy:v1.29-latest \
        --config-path /etc/envoy/envoy.yaml
    ```





