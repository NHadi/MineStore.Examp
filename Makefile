BINARY   := ministore
CMD      := ./cmd/ministore
SQLITE_PKG := modernc.org/sqlite

.PHONY: build test clean cross cross-linux cross-darwin tidy install-deps

# Default: build for current platform
build:
	go build -tags sqlite -o bin/$(BINARY) $(CMD)

# Run all tests
test:
	go test -v ./...

# Cross-platform build: linux/amd64 and darwin/arm64
cross: cross-linux cross-darwin

cross-linux:
	GOOS=linux GOARCH=amd64 go build -tags sqlite -o bin/$(BINARY)-linux-amd64 $(CMD)

cross-darwin:
	GOOS=darwin GOARCH=arm64 go build -tags sqlite -o bin/$(BINARY)-darwin-arm64 $(CMD)

# Dependency management
tidy:
	go mod tidy

# Add SQLite driver (pure Go, no CGO required)
install-deps:
	go get $(SQLITE_PKG)
	go mod tidy

# Cleanup
clean:
	rm -rf bin/
