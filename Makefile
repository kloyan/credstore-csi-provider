default: test

clean: rm -rf *.o

build:
	CGO_ENABLED=0 go build -o dist/ .

test:
	go test ./...

mod:
	@go mod tidy
