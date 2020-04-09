GoContainer=acme_go
RunGo=docker exec -i ${GoContainer}
buildVersion=${shell git rev-parse --short HEAD}
outputDir=${PWD}/dist
testOutputDir=${outputDir}/tests
pkgDir=${PWD}/pkg/...
internalDir=${PWD}/internal/...

# Development
up: create-dev-env
	@make start-containers
	@${RunGo} sh -c 'make local-migration'

down:
	@make stop-containers

start-containers:
	@docker-compose up --build -d postgres golang

stop-containers:
	@docker-compose down -v

create-dev-env:
	@test -e .env || cp .env.example .env

local-migration:
	@go run tools/migration/migrate.go

run:
	@${RunGo} sh -c 'make api-run'


build:
	@docker build -t todoapi:${buildVersion} -f ./docker/golang/Dockerfile .

# CI
ci-test:
	@${RunGo} sh -c 'make api-test'

# API
api-test:
	@make api-lint
	@make api-unit

api-lint:
	@golangci-lint run ${pkgDir} ${internalDir} -v

api-unit:
	@mkdir -p ${testOutputDir}
	@go clean -testcache
	@go test \
        -cover \
        -coverprofile=cp.out \
        -outputdir=${testOutputDir} \
        -race \
        -v \
        -failfast \
        ${pkgDir} \
        ${internalDir}
	@go tool cover -html=${testOutputDir}/cp.out -o ${testOutputDir}/cp.html

api-run:
	@go run cmd/todo/api.go

api-build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
     	    -ldflags='-w -s -extldflags "-static"' \
     	    -a \
     	    -o ${outputDir}/cmd/todo/api ./cmd/todo/api.go
