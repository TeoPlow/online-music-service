module github.com/TeoPlow/online-music-service

go 1.24.2

require (
	github.com/georgysavva/scany v1.2.3
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/jackc/pgx/v5 v5.7.4
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.9.0
	github.com/zhavkk/Auth-protobuf v0.0.4
	go.uber.org/mock v0.5.2
	google.golang.org/grpc v1.72.0
	google.golang.org/protobuf v1.36.6
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

require (
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.8.0
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.0.6 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgtype v1.6.2 // indirect
	github.com/jackc/pgx/v4 v4.10.1
	github.com/jackc/puddle v1.1.3 // indirect
	golang.org/x/crypto v0.33.0
	golang.org/x/text v0.22.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
)

replace github.com/TeoPlow/online-music-service/src/auth_service/pkg/authpb => ./src/auth_service/pkg/authpb

replace github.com/TeoPlow/online-music-service/src/auth_service/internal/repository/mocks => ./src/auth_service/internal/repository/mocks

replace github.com/TeoPlow/online-music-service/src/auth_service/internal/service => ./src/auth_service/internal/service

replace github.com/TeoPlow/online-music-service/src/auth_service/internal/validation => ./src/auth_service/internal/validation

replace github.com/TeoPlow/online-music-service/src/auth_service/internal/pkg/jwt => ./src/auth_service/internal/pkg/jwt

replace github.com/TeoPlow/online-music-service/src/auth_service/internal/logger => ./src/auth_service/internal/logger

replace github.com/TeoPlow/online-music-service/src/auth_service/internal/storage => ./src/auth_service/internal/storage

replace github.com/TeoPlow/online-music-service/src/auth_service/internal/repository/postgres => ./src/auth_service/internal/repository/postgres
