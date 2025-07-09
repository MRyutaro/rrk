#!/bin/sh
set -e

REPO="MRyutaro/rrk"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
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

# Create install directory if it doesn't exist
mkdir -p "$INSTALL_DIR"

# Install without sudo
mv "/tmp/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"

echo "rrk has been installed to ${INSTALL_DIR}/${BINARY_NAME}"

# Add to PATH if not already present
PATH_SETUP_NEEDED=false
case ":$PATH:" in
    *:"$INSTALL_DIR":*)
        ;;
    *)
        echo "Adding '${INSTALL_DIR}' to your PATH..."
        PATH_SETUP_NEEDED=true
        ;;
esac

# Detect shell for setup
SHELL_NAME=$(basename "$SHELL" 2>/dev/null || echo "unknown")
case "$SHELL_NAME" in
    bash|zsh)
        SHELL_CONFIG_FILE="$HOME/.${SHELL_NAME}rc"
        ;;
    *)
        SHELL_NAME="unknown"
        ;;
esac

# Setup PATH if needed
if [ "$PATH_SETUP_NEEDED" = true ] && [ "$SHELL_NAME" != "unknown" ]; then
    echo "export PATH=\"\$PATH:${INSTALL_DIR}\"" >> "$SHELL_CONFIG_FILE"
    echo "✅ Added ${INSTALL_DIR} to PATH in $SHELL_CONFIG_FILE"
    export PATH="$PATH:${INSTALL_DIR}"
fi

# Setup shell integration
echo ""
echo "🔧 Setting up shell integration..."
if [ "$SHELL_NAME" != "unknown" ]; then
    echo "Detected shell: $SHELL_NAME"
    echo "Setting up rrk shell integration..."
    if "${INSTALL_DIR}/${BINARY_NAME}" setup -y 2>/dev/null; then
        echo "✅ Shell integration setup complete!"
        echo ""
        echo "🎉 Installation complete!"
        echo "Please restart your shell or run: source $SHELL_CONFIG_FILE"
    else
        echo "⚠️  Shell integration setup failed. You can set it up later with: rrk setup"
    fi
else
    echo "Could not detect shell. You can set up shell integration later with: rrk setup"
fi

echo ""
echo "Run 'rrk --help' to get started!"
