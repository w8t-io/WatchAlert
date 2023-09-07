FROM registry.js.design/base/golang:1.18 AS build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://goproxy.io"

WORKDIR /root

COPY . /root

RUN go mod tidy && \
    go build -o alertEventMgr ./main.go && \
    chmod 777 alertEventMgr

FROM registry.js.design/base/busybox:glibc

RUN mkdir -p /app/config

COPY --from=build /root/alertEventMgr /app/alertEventMgr

WORKDIR /app

ENTRYPOINT ["/app/alertEventMgr"]