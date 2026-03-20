BINARY_NAME=nostr-relay

# The default 'all' target
all: build

# 1. Tidy up dependencies
tidy:
	go mod tidy

# 2. Build the binary (runs tidy first)
build: tidy
	go build -o $(BINARY_NAME) main.go

# 3. Build and run
run: build
	./$(BINARY_NAME)

# 4. Clean up files
clean:
	rm -f $(BINARY_NAME)
	rm -f relay.bolt

.PHONY: all tidy build run clean
