build:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" ./...

test:
	go test -v

format:
	go fmt ./...
