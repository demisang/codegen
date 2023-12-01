.PHONY: dev-tools fmt tidy lint test
GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vue/*" -not -path "./templates-example/*" | tr "\n" " ")

dev-tools:
	go install github.com/daixiang0/gci@v0.8.0
	go install mvdan.cc/gofumpt@v0.5.0

run:
	go run main.go --root="." --templates="./templates/go-onion" --browser=false

build:
	go build main.go

tidy:
	go mod tidy

fmt:
	gofumpt -w .
	gci write --skip-generated --custom-order -s standard -s default $(GO_FILES)

lint: fmt tidy
	golangci-lint run ./cmd/... ./internal/... ./tests/...

test:
	go test -v ./tests

pack-examples:
	tar czf templates/go_onion.tar.gz -C templates/go-onion .
