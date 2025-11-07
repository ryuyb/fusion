# Docker Multi-Architecture Build Guide

This document explains how to build and deploy Fusion using Docker with support for multiple architectures (AMD64 and ARM64).

## Prerequisites

- Docker Desktop (includes Buildx) or Docker Engine with Buildx plugin
- For multi-architecture builds: `docker buildx` command

Check if Buildx is available:
```bash
docker buildx version
```

## Quick Start

### Build for Multiple Architectures

Use the provided build script (recommended):

```bash
# Build for AMD64 and ARM64
./scripts/build-multiarch.sh

# Build and push to registry
./scripts/build-multiarch.sh --push

# Build for specific platform only
./scripts/build-multiarch.sh --platform linux/arm64

# Build ARM64 and load into local Docker
./scripts/build-multiarch.sh --load --platform linux/arm64
```

### Build Options

The `scripts/build-multiarch.sh` script supports the following options:

- `--push`: Push images to registry after build
- `--load`: Load image into local Docker (single platform only)
- `--platform PLATFORMS`: Specify target platforms (default: `linux/amd64,linux/arm64`)
- `--tag IMAGE_NAME`: Custom image name/tag

### Environment Variables

You can customize the build using environment variables:

```bash
# Set custom image name
IMAGE_NAME=myregistry/fusion ./scripts/build-multiarch.sh

# Set registry prefix
REGISTRY=docker.io/username ./scripts/build-multiarch.sh

# Set version
VERSION=v1.0.0 ./scripts/build-multiarch.sh
```

## Manual Build Commands

### Single Architecture

Build for your current platform:
```bash
docker build -t fusion:latest .
```

Build for specific platform:
```bash
# For ARM64 (Apple Silicon, Raspberry Pi, etc.)
docker build --platform linux/arm64 -t fusion:latest .

# For AMD64 (Intel/AMD processors)
docker build --platform linux/amd64 -t fusion:latest .
```

### Multi-Architecture with Buildx

Create a builder instance (first time only):
```bash
docker buildx create --name fusion-builder --use
docker buildx inspect --bootstrap
```

Build for multiple platforms:
```bash
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  --build-arg VERSION=$(git describe --tags --always --dirty) \
  --build-arg BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S') \
  --build-arg GIT_COMMIT=$(git rev-parse --short HEAD) \
  --build-arg GO_VERSION=$(go version | awk '{print $3}') \
  -t fusion:latest \
  --push \
  .
```

**Note:** Multi-platform builds with `--push` flag require a registry. Use `--load` for local builds, but only with a single platform.

## Docker Compose

The `docker-compose.yml` file will automatically pull the correct image for your platform:

```bash
# Start services
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop services
docker-compose down
```

## Platform-Specific Notes

### ARM64 (Apple Silicon, Raspberry Pi)

The image is optimized for ARM64 and will run natively without emulation:

```bash
# Build locally
./scripts/build-multiarch.sh --load --platform linux/arm64

# Or use docker-compose
docker-compose up
```

### AMD64 (Intel/AMD)

Standard x86_64 build:

```bash
# Build locally
./scripts/build-multiarch.sh --load --platform linux/amd64

# Or use docker-compose
docker-compose up
```

## Image Details

### Build Optimizations

- **Multi-stage build**: Separate build and runtime stages
- **Minimal base image**: Alpine Linux for small image size
- **Static binary**: CGO_ENABLED=0 for portability
- **Build flags**: `-w -s` to strip debug info and reduce size

### Security Features

- Non-root user execution
- Minimal runtime dependencies (only ca-certificates and tzdata)
- Health check endpoint

### Image Size

Approximate sizes:
- AMD64: ~20-25 MB
- ARM64: ~20-25 MB

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Build Multi-Arch Docker Image

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2

      - name: Login to DockerHub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Build and push
        run: |
          chmod +x scripts/build-multiarch.sh
          ./scripts/build-multiarch.sh --push
```

### GitLab CI Example

```yaml
build-multiarch:
  image: docker:latest
  services:
    - docker:dind
  before_script:
    - docker buildx create --use
  script:
    - chmod +x scripts/build-multiarch.sh
    - ./scripts/build-multiarch.sh --push
  only:
    - tags
```

## Troubleshooting

### Buildx Not Available

Install Docker Buildx:
```bash
# For Linux
curl -LO https://github.com/docker/buildx/releases/download/v0.12.0/buildx-v0.12.0.linux-amd64
mkdir -p ~/.docker/cli-plugins
mv buildx-v0.12.0.linux-amd64 ~/.docker/cli-plugins/docker-buildx
chmod +x ~/.docker/cli-plugins/docker-buildx
```

### Builder Instance Issues

Reset builder:
```bash
docker buildx rm fusion-builder
docker buildx create --name fusion-builder --use
docker buildx inspect --bootstrap
```

### Load vs Push

- Use `--load` to load the image into local Docker (single platform only)
- Use `--push` to push to a registry (supports multiple platforms)
- Cannot use both together

### Platform Not Supported

Ensure your Docker installation supports the target platform:
```bash
docker buildx ls
```

## Additional Resources

- [Docker Buildx Documentation](https://docs.docker.com/buildx/working-with-buildx/)
- [Multi-platform Images](https://docs.docker.com/build/building/multi-platform/)
- [Dockerfile Best Practices](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)