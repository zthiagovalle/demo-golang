uid := $(shell id -u)

.PHONY: fmt run test cover mock build update-module doc init install-dependencies env-up env-down migration-up migration-down seed

fmt:
	go fmt ./...

run:
	go run main.go

test: mock
	go-acc --covermode=set -o coverage.txt ./...
	grep -v -E "main.go|_mock.go" coverage.txt > filtered_coverage.txt
	mv filtered_coverage.txt coverage.txt

cover:
	go tool cover -html coverage.txt

mock:
	find . -type f -name "*_mock.go" -exec rm -f {} \;
	go generate -v ./...

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -buildvcs=false -ldflags="-w -s" -o application

update-module:
	go mod tidy

doc:
	swag init --pd --ot=json

init: mock update-module
	cp .env_example .env

install-dependencies:
	go install go.uber.org/mock/mockgen@latest
	go install github.com/ory/go-acc@latest
	go install github.com/swaggo/swag/cmd/swag@latest

env-up:
	cd development-environment && make up

env-down:
	cd development-environment && make down

migration-up:
	docker run --rm --user $(uid) -v ./migrations:/migrations --network host migrate/migrate -database "postgres://demo:demo@localhost:5432/demo?sslmode=disable" -path /migrations up

migration-down:
	docker run --rm --user $(uid) -v ./migrations:/migrations --network host migrate/migrate -database "postgres://demo:demo@localhost:5432/demo?sslmode=disable" -path /migrations down 1

seed:
	docker exec -i demo-postgres psql -U demo -d demo \
	  < development-environment/database/tests-dataset/products-insert.sql
