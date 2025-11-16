# Build stage
FROM --platform=$BUILDPLATFORM golang:1.25.0-alpine AS builder

# Build arguments for cross-compilation
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT=unknown
ARG GO_VERSION

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-w -s \
    -X main.Version=${VERSION} \
    -X main.BuildTime=${BUILD_TIME} \
    -X main.GitCommit=${GIT_COMMIT} \
    -X main.GoVersion=${GO_VERSION}" \
    -o fusion cmd/app/main.go

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -S fusion && adduser -S fusion -G fusion

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/fusion .

# Copy config files
COPY --from=builder /build/configs ./configs

# Change ownership
RUN chown -R fusion:fusion /app

# Switch to non-root user
USER fusion

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

# Run the application
ENTRYPOINT ["./fusion"]
CMD ["serve"]