FROM golang:1.21-alpine3.20 AS build
ARG VERSION

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /root

COPY . /root

RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags="-X main.Version=${VERSION}" -o watchAlert ./main.go && \
    chmod 777 watchAlert

FROM alpine:3.20

COPY --from=build /root/watchAlert /app/watchAlert

WORKDIR /app

ENTRYPOINT ["/app/watchAlert"]