# RivalScope(竞品观察)构建脚本
# 用法:make <target>,常用 make help 查看全部命令

BINARY   := bin/rivalscope
WEB_DIR  := web
CONFIG  ?= config/config.yaml
EMAIL   ?=
PASSWORD ?=
NICKNAME ?=

.PHONY: help dev-api dev-web web build run migrate tidy vet clean

help: ## 查看全部命令
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

dev-api: ## 本地开发:启动后端(127.0.0.1:8000)
	go run . -config $(CONFIG)

dev-web: ## 本地开发:启动前端 Vite(5173,/api 代理到 8000)
	cd $(WEB_DIR) && pnpm dev

web: ## 构建前端产物到 web/dist(go:embed 的嵌入来源)
	cd $(WEB_DIR) && pnpm install && pnpm build

build: web ## 完整构建:前端 + 后端,输出单二进制 $(BINARY)
	go build -o $(BINARY) .

run: build ## 构建并启动服务
	$(BINARY) -config $(CONFIG)

# 数据库迁移,CMD 可取 up/down/status/version,默认 status(登录账号在 config.yaml 的 auth.users 配置)
MIGRATE_CMD ?= status
migrate: ## 数据库迁移:make migrate [CMD=up|down|status|version]
	go run . migrate $(MIGRATE_CMD) -config $(CONFIG)

tidy: ## 整理 Go 依赖
	go mod tidy

vet: ## Go 静态检查
	go vet ./...

clean: ## 清理构建产物
	rm -rf bin $(WEB_DIR)/dist/assets $(WEB_DIR)/dist/index.html
