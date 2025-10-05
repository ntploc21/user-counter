update:
	go mod tidy
	go mod vendor

run:
	go run cmd/main.go

lint:
	gofumpt -l -w ./. && golangci-lint run ./...

test:
	go test ./...
