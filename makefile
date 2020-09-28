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

test-e2e:
	@go test -v ./e2e