.PHONY: test lint build tools setup

LEFTHOOK_VERSION := v2.1.1
GOLANGCI_LINT_VERSION := v2.10.1
XK6_VERSION := v1.4.10

test:
	go test -v -count=1 -timeout 180s ./...

lint:
	golangci-lint run ./...

build:
	mkdir -p bin
	xk6 build --with github.com/ishinvin/xk6-iroha=. --output bin/k6-iroha

tools:
	go install github.com/evilmartians/lefthook/v2@$(LEFTHOOK_VERSION)
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	go install go.k6.io/xk6/cmd/xk6@$(XK6_VERSION)

setup: tools
	lefthook install
