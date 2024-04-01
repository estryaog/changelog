all: build

build:
	@go build -o main cmd/changelog/main.go

run:
	@go run cmd/changelog/main.go

docker-run:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		docker-compose up; \
	fi

docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		docker-compose down; \
	fi