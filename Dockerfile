# 基础镜像
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制所有代码
COPY . .

# 编译Go程序
RUN CGO_ENABLED=0 GOOS=linux go build -o lightchat ./main.go

# 最终镜像
FROM alpine:3.19

# 设置工作目录
WORKDIR /app

# 复制编译后的程序
COPY --from=builder /app/lightchat .
COPY --from=builder /app/emergency ./emergency

# 添加执行权限
RUN chmod +x /app/lightchat
RUN chmod +x /app/emergency/*.sh

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["./lightchat"]