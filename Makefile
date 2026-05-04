BINARY := caseforge

.PHONY: build test vet fmt tidy acceptance install clean help

.DEFAULT_GOAL := help

help:
	@printf "%-15s %s\n" "build"      "Compile the caseforge binary"
	@printf "%-15s %s\n" "test"       "Run all unit tests"
	@printf "%-15s %s\n" "vet"        "Run go vet"
	@printf "%-15s %s\n" "fmt"        "Format source with gofmt"
	@printf "%-15s %s\n" "tidy"       "Run go mod tidy"
	@printf "%-15s %s\n" "acceptance" "Run the full acceptance suite"
	@printf "%-15s %s\n" "install"    "Install binary to GOPATH/bin"
	@printf "%-15s %s\n" "clean"      "Remove built binary"

build:
	go build -o $(BINARY) .

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -l -w .

tidy:
	go mod tidy

acceptance:
	./scripts/acceptance.sh

install:
	go install .

clean:
	rm -f $(BINARY)
