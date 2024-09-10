build:
	go build -o app -v ./...

run: build
	./app

test:
	go test -v ./...