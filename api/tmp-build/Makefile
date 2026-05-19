.PHONY: run build test clean docker-up docker-down migrate fmt

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

test:
	go test ./... -v -cover

test-integration:
	go test ./tests/integration/... -v -tags=integration

clean:
	rm -rf bin/

docker-up:
	docker-compose up -d --build

docker-down:
	docker-compose down

migrate:
	go run cmd/server/main.go --migrate

fmt:
	gofmt -w .
	goimports -w .
