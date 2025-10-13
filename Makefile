# Makefile

# Directories
BIN_DIR := ./bin
CMD_DIR := ./cmd

# Binaries
EXTRACTOR_BIN := $(BIN_DIR)/inkextract
COMICONVERTER_BIN := $(BIN_DIR)/comiconverter

# Commands
EXTRACTOR_CMD := $(CMD_DIR)/inkextract
COMICONVERTER_CMD := $(CMD_DIR)/comiconverter

# Build flags
LDFLAGS := -w -s -extldflags "-static"

.PHONY: all inkextract comiconverter clean

# Default target
all: inkextract comiconverter

# Build inkextract
inkextract:
	CGO_ENABLED=0 go build -ldflags='$(LDFLAGS)' -o $(EXTRACTOR_BIN) $(EXTRACTOR_CMD)

# Build comiconverter
comiconverter:
	CGO_ENABLED=0 go build -ldflags='$(LDFLAGS)' -o $(COMICONVERTER_BIN) $(COMICONVERTER_CMD)

# Clean up build artifacts
clean:
	rm -f $(EXTRACTOR_BIN) $(COMICONVERTER_BIN)

# Install formatters
install-formatters:
	go install github.com/segmentio/golines@latest
	go install mvdan.cc/gofumpt@latest

# Run formatters
run-formatters:
	golines --base-formatter=gofumpt -w .
