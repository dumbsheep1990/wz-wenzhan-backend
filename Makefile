.PHONY: build run clean test deps migrate-up migrate-down dev

# 构建
build:
	go build -o bin/server cmd/server/main.go

# 运行
run: build
	./bin/server

# 开发模式运行
dev:
	go run cmd/server/main.go

# 清理
clean:
	rm -rf bin/
	rm -rf logs/

# 测试
test:
	go test -v ./...

# 安装依赖
deps:
	go mod tidy
	go mod download

# 数据库迁移（需要先创建数据库）
migrate-up:
	@echo "创建数据库表..."
	@echo "请先手动创建数据库: CREATE DATABASE wz_wenzhan CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run

# 创建目录
setup:
	mkdir -p bin logs

# Docker 构建
docker-build:
	docker build -t wz-wenzhan-backend .

# Docker 运行
docker-run:
	docker run -p 8080:8080 wz-wenzhan-backend

# 帮助
help:
	@echo "可用命令:"
	@echo "  build        构建项目"
	@echo "  run          构建并运行项目"
	@echo "  dev          开发模式运行"
	@echo "  clean        清理构建文件"
	@echo "  test         运行测试"
	@echo "  deps         安装依赖"
	@echo "  fmt          格式化代码"
	@echo "  setup        创建必要目录"
	@echo "  help         显示帮助"
