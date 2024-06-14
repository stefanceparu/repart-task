CONTAINER_NAME=reparttask
IMAGE_NAME=reparttask
PORT=8282
IP=$(shell docker exec -it $(CONTAINER_NAME) hostname -i)

test:
	@go clean -testcache && go test ./... -v

build:
	@go build -o bin/reparttask ./cmd/api

run: build
	@bin/reparttask

build-docker:
	@docker build --no-cache --pull -t $(IMAGE_NAME) .
	@docker run -d --name $(CONTAINER_NAME) $(IMAGE_NAME) -p $(PORT):$(PORT) --network=host
	@docker stop $(CONTAINER_NAME)

run-docker:
	@docker start $(CONTAINER_NAME)

stop-docker:
	@docker stop $(CONTAINER_NAME)

clear-docker:
	@docker rm $(CONTAINER_NAME)
	@docker rmi $(CONTAINER_NAME)

get-ip:
	@echo "http://"$(IP):$(PORT)