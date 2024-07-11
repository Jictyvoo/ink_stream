# Makefile

# Directories
BIN_DIR := ./bin
CMD_DIR := ./cmd

# Binaries
EXTRACTOR_BIN := $(BIN_DIR)/extractor
COMICONVERTER_BIN := $(BIN_DIR)/comiconverter

# Commands
EXTRACTOR_CMD := $(CMD_DIR)/extractor
COMICONVERTER_CMD := $(CMD_DIR)/comiconverter

# Build flags
LDFLAGS := -w -s -extldflags "-static"

.PHONY: all extractor comiconverter clean

# Default target
all: extractor comiconverter

# Build extractor
extractor:
	CGO_ENABLED=0 go build -ldflags='$(LDFLAGS)' -o $(EXTRACTOR_BIN) $(EXTRACTOR_CMD)

# Build comiconverter
comiconverter:
	CGO_ENABLED=0 go build -ldflags='$(LDFLAGS)' -o $(COMICONVERTER_BIN) $(COMICONVERTER_CMD)

# Clean up build artifacts
clean:
	rm -f $(EXTRACTOR_BIN) $(COMICONVERTER_BIN)
