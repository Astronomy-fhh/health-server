# 使用 Go 官方镜像
FROM golang:1.21 AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o health-server .

# 使用更小的基础镜像运行
FROM alpine:latest

# 将构建好的二进制文件复制到新的镜像
COPY --from=builder /app/health-server /usr/local/bin/health-server

# 设置默认命令
CMD ["health-server"]
