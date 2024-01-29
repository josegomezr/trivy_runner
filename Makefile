.PHONY: all webapi pullworker scanworker pushworker test format lint

all: webapi pullworker scanworker pushworker

test:
	go test -v ./...

format:
	go fmt ./...
webapi:
	go build -o ./bin/webapi ./cmd/webapi

pullworker:
	go build -o ./bin/pull_worker ./cmd/pullworker

scanworker:
	go build -o ./bin/scan_worker ./cmd/scanworker

pushworker:
	go build -o ./bin/push_worker ./cmd/pushworker

integration-server:
	-docker network create shared_network
	docker-compose -f docker-compose.yml -f docker-compose-integration.yml up

lint:
	docker run --rm -v "$(CURDIR):/app" -w /app golangci/golangci-lint:v1.55.2 golangci-lint run -v


