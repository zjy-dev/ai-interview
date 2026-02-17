GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always 2>/dev/null || echo "dev")

ifeq ($(GOHOSTOS), windows)
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find backend/internal -name '*.proto'")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find backend/api -name '*.proto'")
else
	INTERNAL_PROTO_FILES=$(shell find backend/internal -name '*.proto')
	API_PROTO_FILES=$(shell find backend/api -name '*.proto')
endif

.PHONY: init
# install dependencies and tools
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

.PHONY: third_party
# download third_party proto files
third_party:
	@mkdir -p backend/third_party/google/api
	@mkdir -p backend/third_party/google/protobuf
	@mkdir -p backend/third_party/errors
	@mkdir -p backend/third_party/validate
	curl -sSL -o backend/third_party/google/api/annotations.proto \
		https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto
	curl -sSL -o backend/third_party/google/api/http.proto \
		https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto
	curl -sSL -o backend/third_party/google/api/httpbody.proto \
		https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/httpbody.proto
	curl -sSL -o backend/third_party/google/api/field_behavior.proto \
		https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/field_behavior.proto
	curl -sSL -o backend/third_party/google/api/client.proto \
		https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/client.proto
	curl -sSL -o backend/third_party/google/protobuf/descriptor.proto \
		https://raw.githubusercontent.com/protocolbuffers/protobuf/main/src/google/protobuf/descriptor.proto
	curl -sSL -o backend/third_party/errors/errors.proto \
		https://raw.githubusercontent.com/go-kratos/kratos/main/errors/errors.proto

.PHONY: api
# generate api proto
api:
	protoc --proto_path=./backend/api \
	       --proto_path=./backend/third_party \
	       --go_out=paths=source_relative:./backend/api \
	       --go-http_out=paths=source_relative:./backend/api \
	       --go-grpc_out=paths=source_relative:./backend/api \
	       --openapi_out=fq_schema_naming=true,default_response=false:. \
	       $(API_PROTO_FILES)

.PHONY: config
# generate internal proto (conf)
config:
	protoc --proto_path=./backend/internal \
	       --proto_path=./backend/third_party \
	       --go_out=paths=source_relative:./backend/internal \
	       $(INTERNAL_PROTO_FILES)

.PHONY: sqlc
# generate sqlc code
sqlc:
	cd backend && sqlc generate

.PHONY: wire
# generate wire injection
wire:
	cd backend/cmd/server && wire

.PHONY: generate
# generate all code
generate:
	go generate ./backend/...
	cd backend && go mod tidy

.PHONY: build
# build binary
build:
	cd backend && mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: test
# run unit tests
test:
	cd backend && go test -short -race -coverprofile=coverage.out ./...

.PHONY: test-integration
# run integration tests (requires docker)
test-integration:
	cd backend && go test -race -run Integration ./...

.PHONY: lint
# run linter
lint:
	cd backend && golangci-lint run ./...

.PHONY: migrate-up
# run database migrations up (via sql files)
migrate-up:
	@echo "Apply migrations from backend/sql/migrations/ to your database"
	@echo "Usage: mysql -u root -p ai_interview < backend/sql/migrations/000001_init.up.sql"

.PHONY: migrate-down
# run database migrations down
migrate-down:
	@echo "Rollback: mysql -u root -p ai_interview < backend/sql/migrations/000001_init.down.sql"

.PHONY: dev
# start development server
dev:
	cd backend && go run ./cmd/server/ -conf configs/

.PHONY: docker
# start docker compose services
docker:
	docker compose up -d

.PHONY: all
# generate all then build
all:
	make api
	make config
	make sqlc
	make wire
	make generate
	make build

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
