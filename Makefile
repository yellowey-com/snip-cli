BINARY_NAME=snip
BUILD_DIR=cmd/snip

.PHONY: build install clean

build:
	go build -o $(BINARY_NAME) ./$(BUILD_DIR)

install:
	go build -o $(HOME)/go/bin/$(BINARY_NAME) ./$(BUILD_DIR)

clean:
	rm -f $(BINARY_NAME)
