CC = go
Module = github.com/ChenKS12138/collect-homework-go
Version = `git rev-parse --short HEAD`
Date = `date "+%Y-%m-%d %H:%M:%S"`

install:
	@$(CC) mod download

serve:export DB_DEBUG=true
serve:
	@$(CC) run main.go serve --config-file=.env

migrate:
	@$(CC) run main.go migrate --config-file=.env

migrate-init:
	@$(CC) run main.go migrate --init --config-file=.env

test:export DB_DEBUG=false
test:export EMAIL_PREVENT=true
test:export STORAGE_PATH_PREFIX=../tmp
test:
	@$(CC) test -v ./testing

build: clean
	@env GOOS=darwin GOARCH=amd64 $(CC) build -ldflags "-X '$(Module)/util.Version=$(Version)' -X '$(Module)/util.BuildTime=$(Date)'" -o ./build/main-darwin-64 main.go

build-linux: clean
	@env GOOS=linux GOARCH=amd64 $(CC) build -ldflags "-X '$(Module)/util.Version=$(Version)' -X '$(Module)/util.BuildTime=$(Date)'" -o ./build/main-linux-64 main.go

clean:
	@rm -rf ./build ./tmp

version:
	@echo $(Version)

.PHONY: install serve migrate migrate-init test build build-linux clean version
