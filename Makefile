# Name des Binaries
BINARY_NAME := arch_installer

# Version des Tools (kann manuell gesetzt werden oder aus einer Datei gelesen werden)
VERSION := 1.0.0

# Build-Verzeichnis
BUILD_DIR := build

# Go-Flags
GOFLAGS := -ldflags "-X main.version=$(VERSION)"

# Standardziele
.PHONY: all build test clean version

# Default-Ziel: Alles bauen
all: clean build test

# Build-Ziel: Kompiliere das Tool
build:
	@echo "==> Baue $(BINARY_NAME) Version $(VERSION)"
	@mkdir -p $(BUILD_DIR)
	@go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) main.go
	@echo "==> Binary erstellt: $(BUILD_DIR)/$(BINARY_NAME)"

# Test-Ziel: Führe Tests aus
test:
	@echo "==> Führe Tests aus..."
	@go test ./... -v
	@echo "==> Tests abgeschlossen!"

# Clean-Ziel: Bereinige alte Builds
clean:
	@echo "==> Entferne alte Builds..."
	@rm -rf $(BUILD_DIR)
	@echo "==> Bereinigt."

# Version-Ziel: Zeige die aktuelle Version
version:
	@echo "==> Aktuelle Version: $(VERSION)"
