.DEFAULT_GOAL := build

build:
	@go build -o ./bin/relgen

test:
	@go test -cover ./...

clean:
	@rm -rf ./bin