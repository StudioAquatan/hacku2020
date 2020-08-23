all: build

build: fmt vet
	go build -o bin/oinori

test: fmt vet
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

docker-build:
	docker -t studioaquatan/oinori:$(git rev-parse --short HEAD) .

