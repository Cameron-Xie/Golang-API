DockerGo=docker-compose run --rm --no-deps golang

up:
	@make create-dev-env
	@make install
	@make start-containers

down:
	@make stop-containers

start-containers:
	@docker-compose up --build -d

stop-containers:
	@docker-compose down -v

create-dev-env:
	@test -e .env || cp .env.example .env

install:
	@${DockerGo} go get -d -t -v ./...

run:
	@${DockerGo} go run *.go

test:
	@make acceptance-test

clean-test-cache:
	@${DockerGo} go clean -testcache

acceptance-test:
	@${DockerGo} go test app/handler -v -cover

.PHONY: test