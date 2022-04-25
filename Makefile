clean:
	rm -f ./sampled
.PHONY: clean

compile:
	go build -o sampled ./cmd/sample/...
.PHONY: compile

prepare:
	go mod download
	go mod tidy
.PHONY: prepare

test:
	go test ./...
.PHONY: test

lint:
	@goimports -l -local "github.com/okteto" -w
	@golangci-lint run --out-format github-actions ./...
.PHONY: lint

start:
	./sampled
.PHONY: start

debug:
	dlv debug --headless --listen=:2345 --log --api-version=2
.PHONY: debug

