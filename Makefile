include .env
MIGRATE=migrate -path=migration -database "$(DATABASE_HOST)" -verbose

devtools:
	@echo "Installing devtools"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest
	go install go.uber.org/mock/mockgen@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/daixiang0/gci@v0.11.2
	go get github.com/google/wire/cmd/wire@latest

generate:
	go generate ./...

all_tests:
	go test -v ./internal/v1/http/handler/... ./internal/v1/biz/... -bench=. -cover  -coverprofile=coverage.out -benchmem -cpu=1,2,3,4 -timeout=500ms


bench_tests:
	go test -v ./internal/v1/http/handler/... ./internal/v1/biz/... -bench=. -benchmem -cpu=1,2,3,4 -timeout=500ms

unit_tests:
	go test -v ./internal/v1/http/handler/... ./internal/v1/biz/...

coverage_tests:
	go test -v ./internal/v1/http/handler/... ./internal/v1/biz/... -cover  -coverprofile=coverage.out

fmt:
	gofumpt -l -w .;gci write ./

db-migrate-up:
		$(MIGRATE) up
db-migrate-down:
		$(MIGRATE) down
db-force:
		@read -p  "Which version do you want to force?" VERSION; \
		$(MIGRATE) force $$VERSION

db-goto:
		@read -p  "Which version do you want to migrate?" VERSION; \
		$(MIGRATE) goto $$VERSION

db-drop:
		$(MIGRATE) drop

db-create-migration:
		@read -p  "What is the name of migration?" NAME; \
		${MIGRATE} create -ext sql -seq -dir migration  $$NAME

db-seed:
	go run ./cmd/.  --seed --fakeData


swagger:
	swag init --parseDependency -g ./cmd/boot.go -o ./docs

wire:
	 cd cmd && wire gen && cd ..

build:
	go build -o bin/price ./cmd/.

run:build
	./bin/price

cron:build
	./bin/price --cron