# this is Makefile about melody in golang.
all: test build

run:
	@echo "Run  ..."
	@go run .
	@echo "You can use melody now!"

build:
	@echo "Build  ..."
	@go build .
	@echo "You can use melody now!"


test:
	go generate ./...
	go test -cover -race ./...
