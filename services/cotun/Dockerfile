FROM golang:1.24.0 AS builder
WORKDIR /app
COPY . .

RUN go env -w CGO_ENABLED=0 && \
    go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy,direct

ARG COTUND_VERSION=v1.2.10
RUN go mod tidy 
RUN go build -ldflags="-s -w -X 'github.com/zgsm-ai/cotun/share.BuildVersion=$COTUND_VERSION'" -o cotund server-main/*.go
RUN chmod 755 cotund

FROM alpine:3.21 AS runtime
ENV env prod
ENV TZ Asia/Shanghai
WORKDIR /app
COPY --from=builder /app/cotund /app/cotund
ENTRYPOINT ["/app/cotund"]
