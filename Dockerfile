FROM golang:latest
# 设置工作昌录
WORKDIR /app
# 复制代码到容器中
COPY . .
# 构建应用
# RUN go build -o main
EXPOSE 8089
# 设置启动命令
CMD ["./go-gin-demo"]