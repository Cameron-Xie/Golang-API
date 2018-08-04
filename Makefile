DockerGo=docker exec acme_go

up:
	@make create-dev-env
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
	@${DockerGo} dep ensure

migrate:
	@${DockerGo} go run cmd/migrate/main.go

run:
	@${DockerGo} go run cmd/api/*.go

clean-test-cache:
	@${DockerGo} go clean -testcache

ci-install:
	@make create-dev-env
	@make install

acceptance-test:
	@${DockerGo} go test app/handler -v -cover

.PHONY: test