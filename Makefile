ROOT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

run:
	go run cmd/api/main.go

generate:
	docker run --rm -v $(ROOT_DIR):/src -w /src sqlc/sqlc generate