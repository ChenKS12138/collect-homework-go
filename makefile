include .env.sh


install:
	@go mod download

serve:export DB_DEBUG=true
serve:
	@go run main.go serve
	# @go run main.go serve 2>&1 > out.log > /dev/null

dev:
	@go run main.go

migrate:
	@go run main.go migrate

migrate-init:
	@go run main.go migrate --init

test-e2e:
	@go test -v ./e2e