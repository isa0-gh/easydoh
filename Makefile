# Makefile for easydoh
# - Local builds allowed only on Linux hosts (explicit).
# - Use `make cross-build` on any host to produce linux binaries.
# - Binary will be installed to /bin/easydoh (not a local bin/ directory).
#
# Usage:
#   make                # build for host if linux, otherwise error
#   make cross-build
#   sudo make install

BINARY := easydoh
PKG := ./cmd/easydoh
BIN_DIR ?= /bin
SERVICE_SRC := internal/deploy/easydoh.service
SERVICE_DEST := /etc/systemd/system/$(BINARY).service

HOST_GOOS   := $(shell go env GOOS)
HOST_GOARCH := $(shell go env GOARCH)

.PHONY: all build cross-build install install-binary install-service install-service-only systemd-enable systemd-start systemd-stop clean

all: build

# Build for host if host is linux; otherwise show error and suggest cross-build.
build:
ifeq ($(HOST_GOOS),linux)
	@echo "Building for linux/$(HOST_GOARCH) -> ./$(BINARY)..."
	CGO_ENABLED=0 GOOS=linux GOARCH=$(HOST_GOARCH) go build -trimpath -ldflags="-s -w" -o ./$(BINARY) $(PKG)
	@echo "Built: ./$(BINARY)"
else
	$(error Host OS '$(HOST_GOOS)' not supported for local builds. Use 'make cross-build' to produce linux binaries from any host.)
endif

# Cross-compile for common linux architectures, outputs into project root (no bin/ dir)
cross-build:
	@echo "Cross-building linux/amd64 -> ./$(BINARY)-linux-amd64..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o ./$(BINARY)-linux-amd64 $(PKG)
	@echo "Cross-building linux/arm64 -> ./$(BINARY)-linux-arm64..."
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o ./$(BINARY)-linux-arm64 $(PKG)
	@echo "Built: ./$(BINARY)-linux-amd64 ./$(BINARY)-linux-arm64"

# Install will copy the linux binary to /bin/easydoh and install+enable the service.
# Requires root.
install: install-binary install-service
	@echo "Install complete."

# Copy the binary to BIN_DIR (/bin). If local host built ./easydoh, use it.
# Otherwise fall back to ./easydoh-linux-amd64 produced by cross-build.
install-binary:
	@if [ "$$(id -u)" -ne 0 ]; then echo "install-binary requires root (sudo)"; exit 1; fi
	@if [ -f ./$(BINARY) ]; then \
		echo "Installing ./$(BINARY) -> $(BIN_DIR)/$(BINARY)"; \
		install -Dm0755 ./$(BINARY) $(BIN_DIR)/$(BINARY); \
	else \
		if [ -f ./$(BINARY)-linux-amd64 ]; then \
			echo "Installing ./$(BINARY)-linux-amd64 -> $(BIN_DIR)/$(BINARY)"; \
			install -Dm0755 ./$(BINARY)-linux-amd64 $(BIN_DIR)/$(BINARY); \
		else \
			echo "No linux binary found. Run 'make' on linux or 'make cross-build' to produce linux binaries."; exit 1; \
		fi \
	fi
	@echo "Binary installed to $(BIN_DIR)/$(BINARY)"

# Install systemd unit (requires root). Uses the repo service file.
install-service:
	@if [ "$$(id -u)" -ne 0 ]; then echo "install-service requires root (sudo)"; exit 1; fi
	@if [ ! -f $(SERVICE_SRC) ]; then echo "service file $(SERVICE_SRC) not found"; exit 1; fi
	install -d -m 0755 /etc/systemd/system
	install -m 0644 $(SERVICE_SRC) $(SERVICE_DEST)
	systemctl daemon-reload
	systemctl enable --now $(BINARY).service
	@echo "Systemd unit installed and started."

# Only install the service (if binary already in place)
install-service-only:
	@$(MAKE) install-service

# Systemd management helpers
systemd-enable:
	@if command -v systemctl >/dev/null 2>&1; then \
		systemctl daemon-reload; \
		systemctl enable $(BINARY).service; \
	else \
		echo "systemctl not found"; \
	fi

systemd-start:
	@if command -v systemctl >/dev/null 2>&1; then \
		systemctl start $(BINARY).service; \
	else \
		echo "systemctl not found"; \
	fi

systemd-stop:
	@if command -v systemctl >/dev/null 2>&1; then \
		systemctl stop $(BINARY).service; \
	else \
		echo "systemctl not found"; \
	fi

clean:
	@echo "Cleaning up..."
	rm -f ./$(BINARY) ./$(BINARY)-linux-amd64 ./$(BINARY)-linux-arm64
