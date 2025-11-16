.DEFAULT_GOAL := help

VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GO_VERSION := $(shell go version | awk '{print $$3}')

LDFLAGS := -ldflags "\
	-X main.Version=$(VERSION) \
	-X main.BuildTime=$(BUILD_TIME) \
    -X main.GitCommit=$(GIT_COMMIT) \
    -X main.GoVersion=$(GO_VERSION)"

# colors
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[34m
NC := \033[0m # No Color

.PHONY: help
help: ## 显示帮助信息
	@echo "$(GREEN)Fusion Makefile Commands$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

.PHONY: install-tools
install-tools: ## 安装开发工具
	@echo "$(BLUE)Installing development tools...$(NC)"
	@go install github.com/air-verse/air@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/vektra/mockery/v3@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)Development tools installed$(NC)"

.PHONY: lint
lint: ## 运行 golangci-lint
	@echo "$(BLUE)Running linter...$(NC)"
	@golangci-lint run
	@echo "$(GREEN)Linting completed$(NC)"

.PHONY: generate-ent
generate-ent: ## 生成Ent代码
	@echo "$(GREEN)Generating Ent code...$(NC)"
	@go generate ./internal/infrastructure/database/ent
	@echo "$(GREEN)Code generation completed$(NC)"

.PHONY: ent-new
ent-new: ## 创建新的Ent schema (使用: make ent-new name=User)
	@if [ -z "$(name)" ]; then \
		echo "$(RED)Error: name is required. Usage: make ent-new name=User$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Creating Ent schema: $(name)$(NC)"
	@go run -mod=mod entgo.io/ent/cmd/ent new --target internal/infrastructure/database/schema $(name)

.PHONY: generate-swagger
generate-swagger: ## 生成swagger文档
	@echo "$(GREEN)Generating Swagger documentation...$(NC)"
	@go generate ./docs.go
	@echo "$(GREEN)Swagger documentation generated in ./docs/api/$(NC)"

.PHONY: format-swagger
format-swagger: ## 格式化 swagger 文档
	@swag fmt

.PHONY: generate-mockery
generate-mockery: ## 生成 mockery 文件
	@echo "$(GREEN)Generating mockery files...$(NC)"
	@mockery
	@echo "$(GREEN)Swagger mockery files done$(NC)"

.PHONY: build
build: ## 构建
	@go build $(LDFLAGS) -o bin/fusion cmd/app/main.go

.PHONY: docker-build
docker-build: ## 构建Docker镜像 (单架构)
	@echo "$(BLUE)Building Docker image for current platform...$(NC)"
	@docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg GO_VERSION=$(GO_VERSION) \
		-t fusion:$(VERSION) \
		-t fusion:latest \
		.
	@echo "$(GREEN)Docker image built successfully$(NC)"

.PHONY: docker-build-multiarch
docker-build-multiarch: ## 构建多架构Docker镜像 (AMD64 + ARM64)
	@echo "$(BLUE)Building multi-architecture Docker images...$(NC)"
	@chmod +x scripts/build-multiarch.sh
	@./scripts/build-multiarch.sh
	@echo "$(GREEN)Multi-architecture build completed$(NC)"

.PHONY: docker-build-arm64
docker-build-arm64: ## 构建ARM64架构Docker镜像
	@echo "$(BLUE)Building Docker image for ARM64...$(NC)"
	@chmod +x scripts/build-multiarch.sh
	@./scripts/build-multiarch.sh --load --platform linux/arm64
	@echo "$(GREEN)ARM64 Docker image built successfully$(NC)"

.PHONY: docker-build-amd64
docker-build-amd64: ## 构建AMD64架构Docker镜像
	@echo "$(BLUE)Building Docker image for AMD64...$(NC)"
	@chmod +x scripts/build-multiarch.sh
	@./scripts/build-multiarch.sh --load --platform linux/amd64
	@echo "$(GREEN)AMD64 Docker image built successfully$(NC)"

.PHONY: docker-push
docker-push: ## 构建并推送多架构镜像到仓库
	@echo "$(BLUE)Building and pushing multi-architecture images...$(NC)"
	@chmod +x scripts/build-multiarch.sh
	@./scripts/build-multiarch.sh --push
	@echo "$(GREEN)Images pushed successfully$(NC)"

.PHONY: docker-up
docker-up: ## 启动Docker Compose服务
	@echo "$(BLUE)Starting services with Docker Compose...$(NC)"
	@docker-compose up -d
	@echo "$(GREEN)Services started$(NC)"

.PHONY: docker-down
docker-down: ## 停止Docker Compose服务
	@echo "$(BLUE)Stopping services...$(NC)"
	@docker-compose down
	@echo "$(GREEN)Services stopped$(NC)"

.PHONY: docker-logs
docker-logs: ## 查看Docker Compose日志
	@docker-compose logs -f app
