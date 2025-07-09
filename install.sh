#!/bin/sh
set -e

REPO="MRyutaro/rrk"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
BINARY_NAME="rrk"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

case "$OS" in
    linux|darwin)
        BINARY="rrk-${OS}-${ARCH}"
        ;;
    mingw*|msys*|cygwin*)
        OS="windows"
        BINARY="rrk-${OS}-${ARCH}.exe"
        ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac

# Get latest release URL
DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${BINARY}"

echo "Downloading rrk for ${OS}/${ARCH}..."
echo "URL: ${DOWNLOAD_URL}"

# Download binary
if command -v curl >/dev/null 2>&1; then
    curl -L -o "/tmp/${BINARY_NAME}" "${DOWNLOAD_URL}"
elif command -v wget >/dev/null 2>&1; then
    wget -O "/tmp/${BINARY_NAME}" "${DOWNLOAD_URL}"
else
    echo "Error: curl or wget is required"
    exit 1
fi

# Make executable
chmod +x "/tmp/${BINARY_NAME}"

# Install (may require sudo)
if [ -w "$INSTALL_DIR" ]; then
    mv "/tmp/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
else
    echo "Installing to ${INSTALL_DIR} (requires sudo)..."
    sudo mv "/tmp/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
fi

echo "rrk has been installed to ${INSTALL_DIR}/${BINARY_NAME}"
echo "Run 'rrk' to get started!"