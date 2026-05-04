BINARY := caseforge
INSTALL_DIR := $(shell go env GOPATH)/bin

.PHONY: build test vet fmt acceptance install clean

build:
	go build -o $(BINARY) .

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -l -w .

acceptance: build
	./scripts/acceptance.sh

install:
	go install .

clean:
	rm -f $(BINARY)
