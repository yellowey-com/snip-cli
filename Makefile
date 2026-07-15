BINARY := snip
INSTALL_PATH := $(HOME)/go/bin/$(BINARY)

.PHONY: build install clean

build:
	go build -o $(BINARY) .

install:
	go build -o $(INSTALL_PATH) .

clean:
	rm -f $(BINARY)
