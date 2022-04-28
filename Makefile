clean:
	rm -f ./sampled
.PHONY: clean

compile:
	go build -o sampled ./cmd/sample/...
.PHONY: compile

gen:
	go generate ./...
.PHONY: gen

prepare:
	go mod download
	go mod tidy
.PHONY: prepare

test:
	go test ./...
.PHONY: test

lint:
	@goimports -local "github.com/okteto" -w -l .
	@golangci-lint run --out-format github-actions ./...
.PHONY: lint

start:
	./sampled
.PHONY: start

debug:
	dlv debug --headless --listen=:2345 --log --api-version=2
.PHONY: debug

build-image:
	okteto build -t okteto.dev/okteto-sample .
.PHONY: build-image
