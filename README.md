# Golang REST API

[![Build Status](https://travis-ci.org/Cameron-Xie/Golang-API.svg?branch=master)](https://travis-ci.org/Cameron-Xie/Golang-API)

A TODO list REST API build with Golang.

**Stacks:**
* Golang 1.14
* Postgres 12 (self signed cert with Docker)

**pkgs/tools:**
* go-chi/chi
* jinzhu/gorm
* golang-migrate/migrate
* golangci/golangci-lint

### Endpoints (CRUD)

`GET /tasks`

List all tasks

**Example:**

```Golang
curl --header "Content-Type: application/json" \
  --request GET \
  http://localhost:8080/tasks/
```

`POST /tasks`

Create new task

**Example:**

```Golang
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"name":"project name","description":"project description"}' \
  http://localhost:8080/tasks/
```

`PATCH /tasks/{id} `

Update a task

**Example:**

```Golang
curl --header "Content-Type: application/json" \
  --request PATCH \
  --data '{"name":"first project","description":"first project description"}' \
  http://localhost:8080/tasks/{uuid}
```

`GET /tasks/{id} `

Get a task

**Example:**

```Golang
curl --header "Content-Type: application/json" \
  --request GET \
  http://localhost:8080/tasks/{uuid}
```

`DELETE /tasks/{id} `

Delete a task

**Example:**

```Golang
curl --header "Content-Type: application/json" \
  --request DELETE \
  http://localhost:8080/tasks/{uuid}
```

### Setup Development / Test Environment

**With Docker**

* Make sure you have `Docker` and `Docker Compose` installed.
* Clone the Repository.
* Run `make up` from app root directory. It may take few more minutes for `install packages` and start the container.
* Run `make run` from app root directory to run the API.
* Open `http://127.0.0.1:8080` (default config) in your web browser.

### Run Test

**With Docker**

* Make sure you have installed all the packages (including tests), or you could run `make install` for this.
* Run `make ci-test` from app root directory. it will run all tests for you.


### Build final image

* This project is using Multi-stage to build the final image, Run `make build` from app root directory.

notice: in order to run the api in the final image, need to provide db connection config through env, 
please check `/cmd/todo/api.go`

### Contributing
Feedback is welcome.