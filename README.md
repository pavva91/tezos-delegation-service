# Tezos Delegation Service - Backend Exercise

[![Go](https://github.com/pavva91/tezos-delegation-service/actions/workflows/go.yml/badge.svg)](https://github.com/pavva91/tezos-delegation-service/actions/workflows/go.yml)

[![Go Report Card](https://goreportcard.com/badge/github.com/pavva91/tezos-delegation-service)](https://goreportcard.com/report/github.com/pavva91/tezos-delegation-service)

![Coverage](https://img.shields.io/badge/Coverage-94.3%25-brightgreen)

Build a service that gathers new delegations made on the Tezos protocol and exposes them through a public API.

## Requirements:

- The service will poll the new delegations from this Tzkt API endpoint: https://api.tzkt.io/#operation/Operations_GetDelegations
- The data aggregation service must store the delegation data in a store of your choice.
- The API must read data from that store.
- For each delegation, save the following information: sender's address, timestamp, amount, and block.
- Expose the collected data through a public API at the endpoint `/xtz/delegations`.
  - The expected response format is:
  ```jsx
  {
    "data": [
      {
          "timestamp": "2022-05-05T06:29:14Z",
          "amount": "125896",
          "delegator": "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL",
          "block": "2338084"
      },
      {
          "timestamp": "2021-05-07T14:48:07Z",
          "amount": "9856354",
          "delegator": "KT1JejNYjmQYh8yw95u5kfQDRuxJcaUPjUnf",
          "block": "1461334"
      }
    ],
  }
  ```
  - The sender’s address is the delegator.
  - The delegations must be listed most recent first.
  - The endpoint takes one optional query parameter `year` , which is specified in the format YYYY and will result in the data being filtered for that year only.
- The code must be tested
- There must be a way (Makefile, Docker, …) to run it locally it easily.

## Solution

### Quickstart

#### Run

```bash
docker compose up
```

#### Development Environment (Hot reload inside container)

```bash
docker compose -f docker/dev/docker-compose.yml up
```

#### Config Files

There are 2 config files:

1. For DB: `docker/dev/.env` (copy from: `docker/dev/example.env`)
2. For Go Application: `config/dev-config.yml` (copy from `config/dev-config.yml`)

#### Run PostgreSQL Database (terminal 1 - docker)

1. `cd <project_root>`
2. `cd docker/dev`
3. `docker compose up`

#### Run Go REST API Service (terminal 2 - binary)

##### Go Build and run binary

1. `cd <project_root>`
2. Build: ` go build -o bin/app-amd64-linux main.go`

- Explicit Builds:
  - `GOOS=linux GOARCH=amd64 go build -o bin/app-amd64-linux main.go`
  - `GOOS=darwin GOARCH=amd64 go build -o bin/app-amd64-darwin main.go`

3. Run binary: `SERVER_ENVIRONMENT="dev" bin/app-amd64-linux`

### Try the running application

Try the running application directly on the Swagger API:

- `http://localhost:8080/swagger/index.html#/`

## Linter

The linter will help to write idiomatic go. On project root run:

```bash
golangci-lint run
```

The linter configuration is: `.golangci.yml`

```yaml
linters:
  enable-all: true
```

## Tests

### Run All Tests Run all tests and open results on browser:

- `go test ./... -coverprofile coverage.out && go tool cover -html=coverage.out`

## Run Environments

Setup 3 environments: dev, stage and prod

### Run Dev Environment

#### Run DB (terminal 1)

1. cd docker/dev
2. docker compose up

#### Run Go REST API Service (terminal 3)

1. Create stage config: (`config/dev-config.yml`)
2. `cd <project_root>`
3. `SERVER_ENVIRONMENT="dev" go run main.go`

### Run Stage Environment

#### Run DB (terminal 1)

1. cd docker/stage
2. docker compose up

#### Run Go REST API Service (terminal 3)

1. Create stage config: (`config/stage-config.yml`)
2. `cd <project_root>`
3. `SERVER_ENVIRONMENT="stage" bin/app-amd64-linux`

### Run Production Environment

#### Run DB (terminal 1)

1. cd docker/prod
2. docker compose up

#### Run Go REST API Service (terminal 3)

1. Create prod config: (`config/prod-config.yml`)
2. `cd <project_root>`
3. `SERVER_ENVIRONMENT="prod" bin/app-amd64-linux`

## DB Management

### PostgreSQL Database

#### SQL

Queries:

1. List Delegations by most recent first

```sql
SELECT * FROM delegations ORDER BY "timestamp" DESC
```

2. List delegations by year

```sql
SELECT * FROM delegations WHERE EXTRACT(YEAR FROM timestamp) = 2023
```

### DB Management inside neovim (vim-dadbod)

DB Management inside neovim through dadbod ([tpope/vim-dadbod](https://github.com/tpope/vim-dadbod), [kristijanhusak/vim-dadbod-ui](https://github.com/kristijanhusak/vim-dadbod-ui), [kristijanhusak/vim-dadbod-completion](https://github.com/kristijanhusak/vim-dadbod-completion)):

1. `:DBUI` (\<leader\>db is `:DBUIToggle`)
2. Connection to db (Add connection):
   - `postgres://\<user\>:\<password\>@localhost:\<port\>/\<db_name\>`
3. In case of default values (dev db)
   - `postgres://postgres:localhost@localhost:5432/postgres`

## REST API

Uses [Gin-Gonic](https://gin-gonic.com/docs/)

## ORM

Uses [GORM](https://gorm.io/)

## Format Code

1. `cd ~/go/src/github.com/pavva91/tezos-delegation-service/`
2. `gofmt -l -s -w .`

## Hot Reload (air)

[Go Air](https://github.com/cosmtrek/air) enables hot reloading in go.
Usage of air:

1. Create of config file (.air.toml):
   - `air init`
2. Run app with hot reload (inside project root):
   - `air`
     Instead of:
   - `go run main.go`
   - `go run main.go server_config.go`

## Config YAML

Uses [Clean Env](https://github.com/ilyakaznacheev/cleanenv):

- go get -u github.com/ilyakaznacheev/cleanenv

## Swagger Docs (Swag)

Uses [Swag](https://github.com/swaggo/swag#how-to-use-it-with-gin)

- `go install github.com/swaggo/swag/cmd/swag@latest`
  Note: Go install tries to install the package into $GOBIN, when $GOBIN=/usr/local/go/bin will not work, works with $GOBIN=~/go/bin
  Initialize Swag (on project root):
- `swag init`
  Then get dependencies:
- `go get -u github.com/swaggo/gin-swagger`
- `go get -u github.com/swaggo/files`

Format Swag Comments:

- `swag fmt`

Swagger API:

- `http://localhost:8080/swagger/index.html#/`

## Error Handling

Inspired by:

- https://blog.depa.do/post/gin-validation-errors-handling

### Logging

Zero allocation JSON logger, [zerolog](https://github.com/rs/zerolog):

- `go get -u github.com/rs/zerolog/log`

### Go Validator

Gin uses [Go Validator v10] (https://github.com/go-playground/validator):

- `go get github.com/go-playground/validator/v10`

In code import:

- `import "github.com/go-playground/validator/v10" `

## Unit Tests

- `go get -u github.com/stretchr/testify`
- Run all tests: `go test ./...`
- Run all tests and create code coverage report: `go test -v -coverprofile cover.out ./...`
- Run all tests with code coverage and open on browser: `go test -v -coverprofile cover.out ./... && go tool cover -html=cover.out`
- Run tests of /controllers with code coverage and open on browser: `go test -v -coverprofile cover.out ./controllers/ && go tool cover -html=cover.out`
- Run tests of /services with code coverage and open on browser: `go test -v -coverprofile cover.out ./services/ && go tool cover -html=cover.out`

### Code Coverage

- By package name:

  - Just run: `go test -cover github.com/pavva91/tezos-delegation-service/validation`
  - Create coverage file: `go test -v -coverprofile cover.out github.com/pavva91/tezos-delegation-service/validation`
  - One Command create coverage file and open in browser; `go test -v -coverprofile cover.out github.com/pavva91/tezos-delegation-service/controllers && go tool cover -html=cover.out`

- By folder:

  - Just run: `go test -cover ./validation`
  - Create coverage file: `go test -v -coverprofile cover.out ./validation`
  - Open coverage file on browser: `go tool cover -html=cover.out`
  - One Command create coverage file and open in browser; `go test -v -coverprofile cover.out ./controllers/ && go tool cover -html=cover.out`

- Run all tests:
  - Just run: `go test ./... -cover`
  - Create coverage file: `go test ./... -coverprofile coverage.out`
  - Open coverage file on browser: `go tool cover -html=coverage.out`

From [stack overflow](https://stackoverflow.com/questions/10516662/how-to-measure-test-coverage-in-go)

1. Create function in ~/.bashrc and/or ~/.zshrc:

```bash
cover () {
  t="/tmp/go-cover.$$.tmp"
  go test -coverprofile=$t $@ && go tool cover -html=$t && unlink $t
}
```

2. Call this function:

- `cd ~/go/src/github.com/pavva91/tezos-delegation-service/ `
- `cover github.com/pavva91/tezos-delegation-service/validation`

### Run Unit Test with Debugger (neovim DAP)

- `:lua require('dap-go').debug_test()`
- Keymap: \<leader\>dt

## DB Management

### DB Management inside neovim (vim-dadbod)

DB Management inside neovim through dadbod ([tpope/vim-dadbod](https://github.com/tpope/vim-dadbod), [kristijanhusak/vim-dadbod-ui](https://github.com/kristijanhusak/vim-dadbod-ui), [kristijanhusak/vim-dadbod-completion](https://github.com/kristijanhusak/vim-dadbod-completion)):

1. `:DBUI` (\<leader\>db is `:DBUIToggle`)
2. Connection to db (Add connection):
   - `postgres://\<user\>:\<password\>@localhost:\<port\>/\<db_name\>`
3. In case of default values (dev db)
   - `postgres://postgres:localhost@localhost:5432/postgres`

## Check vulnerabilities

```bash
govulncheck ./...
```

## Run Local Go Doc

```bash
go install golang.org/x/tools/cmd/godoc@latest
godoc -http=:6060
```
