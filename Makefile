BINARY_NAME=nostr-relay
INSTALL_DIR=$(HOME)/go/bin

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

# 4. Install the binary to ~/go/bin
install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY_NAME) $(INSTALL_DIR)
	@echo "Successfully installed $(BINARY_NAME) to $(INSTALL_DIR)"

# 5. Clean up local files
clean:
	rm -f $(BINARY_NAME)
	rm -f relay.bolt

.PHONY: all tidy build run clean install
