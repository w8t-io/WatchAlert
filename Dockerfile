FROM registry.cn-hangzhou.aliyuncs.com/opsre/golang:1.21.9-alpine3.19 AS build

ARG VERSION

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /root

COPY . /root

RUN sed -i "s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g" /etc/apk/repositories \
    && apk upgrade && apk add --no-cache --virtual .build-deps \
    ca-certificates upx

RUN go build --ldflags="-X main.Version=${VERSION}" -o watchAlert . \
  && upx -9 watchAlert && chmod +x watchAlert

FROM registry.cn-hangzhou.aliyuncs.com/opsre/alpine:3.19

COPY --from=build /root/watchAlert /app/watchAlert

WORKDIR /app

ENTRYPOINT ["/app/watchAlert"]