############################################################
# Dockerfile to build golang Installed Containers

# Based on alpine

############################################################

FROM golang:1.22.3 AS builder

COPY . /src
WORKDIR /src

RUN GOPROXY="https://goproxy.cn,direct" CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o zelda cmd/server/main.go

FROM alpine:3.13

RUN mkdir /zelda
COPY --from=builder /src/zelda /zelda

EXPOSE 8080
WORKDIR /zelda
CMD ["./zelda"]
