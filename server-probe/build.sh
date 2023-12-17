#!/bin/bash

# 保存原始的 GOOS 和 GOARCH 环境变量的值
ORIGINAL_GOOS="$GOOS"
ORIGINAL_GOARCH="$GOARCH"

# 定义目标平台和架构
TARGET_OS="linux"
TARGET_ARCH=("amd64")

# 清理之前的构建结果
rm -rf build

# 创建构建目录
mkdir -p build

# 遍历目标架构，进行交叉编译
for arch in "${TARGET_ARCH[@]}"; do
    output_name="hello_${arch}"

    if [ "$arch" == "arm" ]; then
        output_name+="v6"
    fi

    if [ "$arch" == "arm64" ]; then
        output_name+="v8"
    fi

    # 执行交叉编译
    GOOS="$TARGET_OS" GOARCH="$arch" 
    go build -o "build/$output_name" center/main.go
    go build -o "build/$output_name" probe/main.go

    if [ $? -eq 0 ]; then
        echo "Successfully built for $TARGET_OS/$arch"
    else
        echo "Failed to build for $TARGET_OS/$arch"
    fi
done

# 恢复原来的配置
export GOOS="$ORIGINAL_GOOS"
export GOARCH="$ORIGINAL_GOARCH"