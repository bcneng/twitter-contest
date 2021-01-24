ldflags = -X main.Version=$(shell git rev-parse HEAD)

build:
	GOOS=linux GOARCH=amd64 go build -ldflags "$(ldflags)" -o functions/contest main.go