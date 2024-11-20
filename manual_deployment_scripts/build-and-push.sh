#!/bin/bash

# Get the script's directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
# Get the project root directory (one level up from script directory)
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"

# Configuration
S3_BUCKET="sns-tool-binaries"
AWS_PROFILE="${AWS_PROFILE:-my}"  # Use specified profile or default

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Building binaries...${NC}"
echo -e "Project root: ${PROJECT_ROOT}"

# Change to project root directory
cd "$PROJECT_ROOT" || exit 1

# Create manual_dist directory
mkdir -p manual_dist

# Build for each platform
PLATFORMS=("linux/amd64" "darwin/amd64" "windows/amd64")

for platform in "${PLATFORMS[@]}"; do
    # Split platform into OS and arch
    IFS="/" read -r -a platform_split <<< "$platform"
    GOOS="${platform_split[0]}"
    GOARCH="${platform_split[1]}"
    
    # Set binary name based on OS
    if [ "$GOOS" = "windows" ]; then
        binary_name="sns-tool.exe"
    else
        binary_name="sns-tool"
    fi
    
    output_name="manual_dist/sns-tool-${GOOS}-${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        output_name="${output_name}.exe"
    fi

    echo -e "${GREEN}Building for ${GOOS}/${GOARCH}...${NC}"
    
    # Build the binary
    GOOS=$GOOS GOARCH=$GOARCH go build -o "$output_name" .
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}Successfully built for ${GOOS}/${GOARCH}${NC}"
    else
        echo -e "${RED}Failed to build for ${GOOS}/${GOARCH}${NC}"
        exit 1
    fi
done

echo -e "${GREEN}Uploading to S3...${NC}"

# Upload to S3 with proper content type and ACL
for file in manual_dist/*; do
    filename=$(basename "$file")
    
    # Set content type based on file extension
    if [[ "$filename" == *.exe ]]; then
        content_type="application/vnd.microsoft.portable-executable"
    else
        content_type="application/octet-stream"
    fi
    
    echo "Uploading $filename to S3..."
    aws s3 cp "$file" "s3://${S3_BUCKET}/$filename" \
        --content-type "$content_type" \
        --acl public-read \
        --profile "$AWS_PROFILE"
        
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}Successfully uploaded $filename${NC}"
    else
        echo -e "${RED}Failed to upload $filename${NC}"
        exit 1
    fi
done

# Upload install script
echo -e "${GREEN}Uploading install script...${NC}"
aws s3 cp "${SCRIPT_DIR}/install.sh" "s3://${S3_BUCKET}/install.sh" \
    --content-type "text/x-shellscript" \
    --acl public-read \
    --profile "$AWS_PROFILE"

echo -e "${GREEN}Build and upload complete!${NC}"
echo -e "Install script URL: https://${S3_BUCKET}.s3.amazonaws.com/install.sh"