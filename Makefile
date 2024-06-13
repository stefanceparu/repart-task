test:
	@go clean -testcache && go test ./... -v

build:
	@go build -o bin/reparttask ./cmd/api

run: build
	@bin/reparttask