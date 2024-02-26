# 使用官方 Gin 镜像作为基础镜像
FROM ubuntu:latest

# 将工作目录设置为 /app
WORKDIR /app

# 将当前目录的内容复制到容器的 /app 中
COPY . /app

RUN apt update && apt upgrade
RUN apt golang-go

## 构建应用
#RUN go build -o coolrun_app_docker_server

# 设置容器端口映射
EXPOSE 8080

# 定义容器启动命令
CMD ["go run /app/main.go"]