install:
	@go mod download

serve:export DB_DEBUG=true
serve:
	@go run main.go serve --config-file=.env

dev:
	@go run main.go

migrate:
	@go run main.go migrate --config-file=.env

migrate-init:
	@go run main.go migrate --init --config-file=.env

test:export EMAIL_PREVENT=true
test:export STORAGE_PATH_PREFIX=../tmp
test:
	@go test -v ./testing

build:
	@env GOOS=darwin GOARCH=amd64 go build -o ./build/main-darwin-64 main.go

build-linux:
	@env GOOS=linux GOARCH=amd64 go build -o ./build/main-linux-64 main.go

clean:
	@rm -rf ./build ./tmp