############################################################
# Dockerfile to build golang Installed Containers

# Based on alpine

############################################################

FROM golang:1.22.3 AS builder

COPY . /src
WORKDIR /src

RUN GOPROXY="https://goproxy.cn,direct" CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o core cmd/server/main.go

FROM alpine:3.13

RUN mkdir /core
COPY --from=builder /src/core /core

EXPOSE 8080
WORKDIR /core
CMD ["./core"]
