FROM golang:1.25-alpine AS builder

RUN go env -w GO111MODULE=auto \
  && go env -w CGO_ENABLED=0 \
  && go env -w GOPROXY=https://goproxy.cn,direct

WORKDIR /build
COPY ./ .
RUN set -ex \
  && cd /build \
  && go mod tidy \
  && go build -ldflags "-s -w -extldflags '-static'" -o xrcuo-api



