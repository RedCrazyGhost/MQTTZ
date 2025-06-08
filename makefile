# 项目基本信息
AUTHOR := "RedCrazyGhost"
REPO_URL := "https://github.com/RedCrazyGhost/MQTTZ"
BINARY_NAME := MQTTZ
CONFIG_PATH := ./bin/conf/config.yaml

# 工具检查
GOLANGCI_LINT := $(shell command -v golangci-lint 2> /dev/null)
WIRE := $(shell command -v wire 2> /dev/null)

# 版本信息
CURRENT_TAG := $(strip $(shell git describe --tags --always))
CURRENT_COMMIT_HASH := $(strip $(shell git rev-parse --short HEAD))

# 版本号处理
ifeq ($(CURRENT_TAG),$(CURRENT_COMMIT_HASH))
  VERSION := $(CURRENT_COMMIT_HASH)
else
  ifeq ($(CURRENT_COMMIT_HASH),$(shell git show-ref --hash=7 $(CURRENT_TAG)))
    VERSION := $(shell git describe --tags --abbrev=0)
  else
    VERSION := $(shell git describe --tags --abbrev=0)-$(CURRENT_COMMIT_HASH)
  endif
endif

# 构建参数
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
LDFLAGS := -w -s -X main.Version=$(VERSION) -X main.Author=$(AUTHOR) -X main.RepoURL=$(REPO_URL)

# 默认目标
.PHONY: all
all: clean lint build

# 清理构建产物
.PHONY: clean
clean:
	@echo "🧹 清理构建产物..."
	@rm -rf ./bin/$(BINARY_NAME)
	@echo "✅ 清理完成"

# 代码检查
.PHONY: lint
lint:
ifndef GOLANGCI_LINT
	@echo "🚀 golangci-lint 未安装"
	@echo "📦 正在安装 golangci-lint v2..."
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6
	@echo "✅ golangci-lint v2 安装完成"
else
	@echo "✅ golangci-lint 已安装"
	@echo "🔄 检查 golangci-lint 版本..."
	@golangci-lint --version | grep -q "v2" || (echo "⚠️ 当前版本不是 v2，正在更新..." && go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6)
endif
	@echo "🔍 开始代码检查..."
	@golangci-lint run
	@echo "✨ 代码检查完成"

# 本地构建
.PHONY: build
build: wire
	@echo "🔨 开始构建..."
	@go fmt ./...
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ./bin/$(BINARY_NAME) -ldflags "$(LDFLAGS)" ./cmd
	@echo "✅ 构建完成"

# 运行应用
.PHONY: run
run: build
	@echo "🚀 启动应用..."
	@./bin/$(BINARY_NAME) -config $(CONFIG_PATH)

# Wire 依赖注入
.PHONY: wire
wire:
ifndef WIRE
	@echo "🚀 wire 未安装"
	@echo "📦 正在安装 wire..."
	@go install github.com/google/wire/cmd/wire@latest
	@echo "✅ wire 安装完成"
else
	@echo "✅ wire 已安装"
endif
	@echo "🔧 生成依赖注入代码..."
	@wire ./cmd
	@echo "✅ 依赖注入代码生成完成"

# 显示版本信息
.PHONY: version
version:
	@echo "📦 当前版本: $(VERSION)"

# 帮助信息
.PHONY: help
help:
	@echo "📋 可用命令:"
	@echo "  make all      - 清理、检查并构建项目"
	@echo "  make clean    - 清理构建产物"
	@echo "  make lint     - 运行代码检查"
	@echo "  make build    - 构建项目"
	@echo "  make run      - 构建并运行项目"
	@echo "  make wire     - 生成依赖注入代码"
	@echo "  make version  - 显示当前版本信息"