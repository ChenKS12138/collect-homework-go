CC = go

install:
	@$(CC) mod download

serve:export DB_DEBUG=true
serve:
	@$(CC) run main.go serve --config-file=.env

migrate:
	@$(CC) run main.go migrate --config-file=.env

migrate-init:
	@$(CC) run main.go migrate --init --config-file=.env

test:export EMAIL_PREVENT=true
test:export STORAGE_PATH_PREFIX=../tmp
test:
	@$(CC) test -v ./testing

build: clean
	@env GOOS=darwin GOARCH=amd64 $(CC) build -o ./build/main-darwin-64 main.go

build-linux: clean
	@env GOOS=linux GOARCH=amd64 $(CC) build -o ./build/main-linux-64 main.go

clean:
	@rm -rf ./build ./tmp

.PHONY: install serve migrate migrate-init test build build-linux clean