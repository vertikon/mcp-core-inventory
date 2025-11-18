.PHONY: run test lint build telemetry

run:
	GO_ENV=dev go run ./cmd

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/{{ .ServiceName }} ./cmd

lint:
	golangci-lint run ./...

test:
	go test ./...

telemetry:
	opentelemetry-collector --config ./configs/otel.yaml
