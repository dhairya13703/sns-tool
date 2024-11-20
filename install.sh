#!/bin/bash

# Configuration
S3_BUCKET="sns-tool-binaries"
S3_BASE_URL="https://${S3_BUCKET}.s3.amazonaws.com"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Determine OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Convert architecture names
case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

# Handle OS detection
case $OS in
    linux|darwin)
        ;;
    *)
        echo -e "${RED}Unsupported operating system: $OS${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}Downloading sns-tool for $OS ($ARCH)...${NC}"

# Construct binary name and URL
BINARY_NAME="sns-tool-${OS}-${ARCH}"
if [ "$OS" = "windows" ]; then
    BINARY_NAME="${BINARY_NAME}.exe"
fi
BINARY_URL="${S3_BASE_URL}/${BINARY_NAME}"

# Download the binary to current directory
echo -e "${GREEN}Downloading from $BINARY_URL...${NC}"
if ! curl -L -f "$BINARY_URL" -o "sns-tool"; then
    echo -e "${RED}Failed to download sns-tool${NC}"
    exit 1
fi

# Make binary executable
chmod +x sns-tool

echo -e "${GREEN}Download complete! sns-tool is in the current directory${NC}"
echo -e "Run ${GREEN}./sns-tool --help${NC} to see available commands"