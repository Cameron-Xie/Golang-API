# Golang REST API

[![Build Status](https://travis-ci.org/Cameron-Xie/Golang-API.svg?branch=master)](https://travis-ci.org/Cameron-Xie/Golang-API)

A Golang API (Starter).

**Stacks:**
* Language: Golang 1.10+
* Containerisation: Docker CE

### Endpoints

`GET /projects/`

List all projects

**Example:**

```Golang
curl --header "Content-Type: application/json" \
  --request GET \
  http://localhost:8080/projects/
```

`POST /projects/`

Create new project

**Example:**

```Golang
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"name":"project name","description":"project description"}' \
  http://localhost:8080/projects/
```

`PATCH /projects/{id} `

Update project

**Example:**

```Golang
curl --header "Content-Type: application/json" \
  --request PATCH \
  --data '{"name":"first project","description":"first project description"}' \
  http://localhost:8080/projects/1
```

`DELETE /projects/{id} `

Delete project

**Example:**

```Golang
curl --header "Content-Type: application/json" \
  --request DELETE \
  http://localhost:8080/projects/1
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
* Run `make test` from app root directory. it will run all tests for you.

### Contributing
Feedback is welcome.