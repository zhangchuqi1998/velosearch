.PHONY: build test lint proto bench docker clean tidy

build:
	go build -o bin/velosearch.exe ./cmd/server

test:
	go test -race -count=1 ./...

lint:
	golangci-lint run ./...

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       proto/velosearch.proto

bench:
	go test -bench=. -benchmem -run=^$$ ./...

docker:
	docker build -t velosearch:dev -f deploy/Dockerfile .

tidy:
	go mod tidy

clean:
	if exist bin rmdir /s /q bin
