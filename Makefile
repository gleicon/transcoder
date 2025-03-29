# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod
BINARY_NAME=transcoder

# Build flags
LDFLAGS=-ldflags "-w -s"

# Test flags
TESTFLAGS=-v -race -coverprofile=coverage.out

# OS and architecture detection
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)
ifeq ($(UNAME_S),Darwin)
	OS := darwin
	PKG_MANAGER := brew
	INSTALL_CMD := brew install
	UPDATE_CMD := brew update
	# Check if running under Rosetta 2
	IS_ROSETTA := $(shell sysctl -n sysctl.proc_translated 2>/dev/null || echo "0")
	ifeq ($(IS_ROSETTA),1)
		ARCH := arm64
		BREW_PREFIX := /opt/homebrew
		BREW_CMD := arch -arm64 brew
	else
		ifeq ($(UNAME_M),arm64)
			ARCH := arm64
			BREW_PREFIX := /opt/homebrew
			BREW_CMD := arch -arm64 brew
		else
			ARCH := x86_64
			BREW_PREFIX := /usr/local
			BREW_CMD := brew
		endif
	endif
else
	OS := linux
	PKG_MANAGER := apt
	INSTALL_CMD := sudo apt install -y
	UPDATE_CMD := sudo apt update
endif

# Default target
all: setup build model

model:
	mkdir models
	curl -L https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.en.bin -o models/ggml-base.en.bin
# Setup development environment
setup:
	@echo "Setting up development environment..."
ifeq ($(OS),darwin)
	@if ! command -v brew &> /dev/null; then \
		/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"; \
		if [ "$(ARCH)" = "arm64" ]; then \
			echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zshrc; \
			eval "$(/opt/homebrew/bin/brew shellenv)"; \
		fi; \
	fi
	$(BREW_CMD) update
	$(BREW_CMD) install go ffmpeg git make cmake
	$(BREW_CMD) install whisper-cpp
else
	$(UPDATE_CMD)
	$(INSTALL_CMD) golang-go ffmpeg libavcodec-extra libavformat-dev libavutil-dev \
		libavdevice-dev libavfilter-dev libswscale-dev libavresample-dev git make curl \
		build-essential cmake
endif
	@echo "Installing Go dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Build the application
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) ./cmd/transcoder

# Run tests
test:
	$(GOTEST) $(TESTFLAGS) ./...

# Run tests with coverage report
test-coverage:
	$(GOTEST) $(TESTFLAGS) ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out
	rm -f coverage.html

# Download dependencies
deps:
	$(GOMOD) download

# Tidy dependencies
tidy:
	$(GOMOD) tidy

# Run the application
run: build
	./$(BINARY_NAME)

# Install the application
install: build
	install -m 755 $(BINARY_NAME) /usr/local/bin/

# Uninstall the application
uninstall:
	rm -f /usr/local/bin/$(BINARY_NAME)

.PHONY: all setup build test test-coverage clean deps tidy run install uninstall 
