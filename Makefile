APP_NAME=go-simple-redis

build:
	go build -o bin/ ./cmd/...

run: build
	./bin/$(APP_NAME)