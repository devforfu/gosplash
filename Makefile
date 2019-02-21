GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
SOURCE_DIR=./src
BINARY_DIR=./bin

sources := $(wildcard $(SOURCE_DIR)/*.go)
executables := $(patsubst $(SOURCE_DIR)/%.go, $(BINARY_DIR)/%, $(sources))

.PHONY: all

all: test build

$(BINARY_DIR)/%: $(SOURCE_DIR)/%.go
	$(GOBUILD) -o $@ $<

build: $(executables)

clean:
	rm -rf $(BINARY_DIR)/**

test:
	$(GOTEST) -v ./src/fetch
	$(GOTEST) -v ./src/links
