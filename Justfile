
run-docker:
    @docker compose up --build   

build:
    go build -o bin/main ./cmd

run:
    go run ./cmd/

test:
    go test -v ./...