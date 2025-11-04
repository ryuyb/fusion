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
	@swag init -g docs.go --output ./docs
	@echo "$(GREEN)Swagger documentation generated in ./docs/$(NC)"

.PHONY: build
build: ## 构建
	@go build $(LDFLAGS) -o bin/fusion cmd/server/main.go
