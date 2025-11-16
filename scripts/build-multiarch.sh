#!/bin/bash

# Multi-architecture Docker build script for Fusion
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Build variables
VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GO_VERSION=$(go version | awk '{print $3}' 2>/dev/null || echo "unknown")
IMAGE_NAME=${IMAGE_NAME:-fusion}
REGISTRY=${REGISTRY:-}

# Add registry prefix if provided
if [ -n "$REGISTRY" ]; then
    IMAGE_NAME="${REGISTRY}/${IMAGE_NAME}"
fi

echo -e "${BLUE}Building multi-architecture Docker images...${NC}"
echo -e "${YELLOW}Version: ${VERSION}${NC}"
echo -e "${YELLOW}Build Time: ${BUILD_TIME}${NC}"
echo -e "${YELLOW}Git Commit: ${GIT_COMMIT}${NC}"
echo -e "${YELLOW}Go Version: ${GO_VERSION}${NC}"
echo -e "${YELLOW}Image Name: ${IMAGE_NAME}${NC}"
echo ""

# Check if buildx is available
if ! docker buildx version &> /dev/null; then
    echo -e "${RED}Error: Docker Buildx is not available${NC}"
    echo -e "${YELLOW}Please install Docker Buildx or use Docker Desktop${NC}"
    exit 1
fi

# Create builder instance if it doesn't exist
if ! docker buildx inspect fusion-builder &> /dev/null; then
    echo -e "${BLUE}Creating buildx builder instance...${NC}"
    docker buildx create --name fusion-builder --use
else
    echo -e "${BLUE}Using existing buildx builder instance...${NC}"
    docker buildx use fusion-builder
fi

# Bootstrap the builder
docker buildx inspect --bootstrap

# Parse command line arguments
PUSH=false
PLATFORMS="linux/amd64,linux/arm64"
LOAD=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --push)
            PUSH=true
            shift
            ;;
        --load)
            LOAD=true
            shift
            ;;
        --platform)
            PLATFORMS="$2"
            shift 2
            ;;
        --tag)
            IMAGE_NAME="$2"
            shift 2
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            echo "Usage: $0 [--push] [--load] [--platform PLATFORMS] [--tag IMAGE_NAME]"
            echo "  --push: Push images to registry"
            echo "  --load: Load image into local docker (single platform only)"
            echo "  --platform: Target platforms (default: linux/amd64,linux/arm64)"
            echo "  --tag: Image name/tag"
            exit 1
            ;;
    esac
done

# Build command
BUILD_CMD="docker buildx build"
BUILD_CMD="$BUILD_CMD --platform ${PLATFORMS}"
BUILD_CMD="$BUILD_CMD --build-arg VERSION=${VERSION}"
BUILD_CMD="$BUILD_CMD --build-arg BUILD_TIME=${BUILD_TIME}"
BUILD_CMD="$BUILD_CMD --build-arg GIT_COMMIT=${GIT_COMMIT}"
BUILD_CMD="$BUILD_CMD --build-arg GO_VERSION=${GO_VERSION}"
BUILD_CMD="$BUILD_CMD -t ${IMAGE_NAME}:${VERSION}"
BUILD_CMD="$BUILD_CMD -t ${IMAGE_NAME}:latest"

if [ "$PUSH" = true ]; then
    BUILD_CMD="$BUILD_CMD --push"
    echo -e "${YELLOW}Images will be pushed to registry${NC}"
elif [ "$LOAD" = true ]; then
    # --load only supports single platform
    if [[ "$PLATFORMS" == *","* ]]; then
        echo -e "${RED}Error: --load can only be used with a single platform${NC}"
        echo -e "${YELLOW}Please specify a single platform with --platform${NC}"
        exit 1
    fi
    BUILD_CMD="$BUILD_CMD --load"
    echo -e "${YELLOW}Image will be loaded into local docker${NC}"
fi

BUILD_CMD="$BUILD_CMD ."

# Execute build
echo -e "${BLUE}Executing build command...${NC}"
echo -e "${YELLOW}$BUILD_CMD${NC}"
echo ""

eval $BUILD_CMD

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✓ Build completed successfully!${NC}"
    echo -e "${GREEN}Images built for platforms: ${PLATFORMS}${NC}"
    echo -e "${GREEN}Tags:${NC}"
    echo -e "  - ${IMAGE_NAME}:${VERSION}"
    echo -e "  - ${IMAGE_NAME}:latest"
else
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi