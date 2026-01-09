# Variables
BINARY_NAME = easydoh
BIN_DIR = /bin
SYSTEMD_SERVICE_SRC = deploy/easydoh.service
SYSTEMD_SERVICE_DEST = /etc/systemd/system/$(BINARY_NAME).service

# Default target
all: build

# Build the Go binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) .

# Install binary and systemd service
install: build
	@echo "Installing $(BINARY_NAME) to $(BIN_DIR)..."
	install -Dm755 $(BINARY_NAME) $(BIN_DIR)/$(BINARY_NAME)
	
	# Install systemd service if systemd is detected
	@if command -v systemctl >/dev/null 2>&1 && [ -d /run/systemd/system ]; then \
		echo "Systemd detected, installing systemd service..."; \
		install -Dm644 $(SYSTEMD_SERVICE_SRC) $(SYSTEMD_SERVICE_DEST); \
	else \
		echo "Systemd not detected. Skipping service installation."; \
	fi
	
	@echo "Installation complete."

# Clean up the built binary
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)

# Systemd service management helpers
systemd-enable:
	@if command -v systemctl >/dev/null 2>&1; then \
		systemctl daemon-reload; \
		systemctl enable $(BINARY_NAME).service; \
	else \
		echo "Systemd not detected."; \
	fi

systemd-start:
	@if command -v systemctl >/dev/null 2>&1; then \
		systemctl start $(BINARY_NAME).service; \
	else \
		echo "Systemd not detected."; \
	fi

systemd-stop:
	@if command -v systemctl >/dev/null 2>&1; then \
		systemctl stop $(BINARY_NAME).service; \
	else \
		echo "Systemd not detected."; \
	fi

.PHONY: all build install clean systemd-enable systemd-start systemd-stop
