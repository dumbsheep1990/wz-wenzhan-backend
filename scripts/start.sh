#!/bin/bash

# 万知文站后端启动脚本

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "Go环境未安装，请先安装Go 1.21+"
    exit 1
fi

# 检查依赖
echo "检查并安装依赖..."
go mod tidy

# 创建必要的目录
mkdir -p logs bin

# 构建应用
echo "构建应用..."
go build -o bin/server cmd/server/main.go

if [ $? -eq 0 ]; then
    echo "构建成功！"
    echo "启动服务..."
    ./bin/server
else
    echo "构建失败！"
    exit 1
fi
