FROM golang:1.24-alpine as builder

RUN sed -i -e 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk --no-cache add git ca-certificates gcc g++

RUN apk add --update gcc musl-dev
RUN apk add --no-cache git
RUN apk add --no-cache sqlite-libs sqlite-dev
RUN apk add --no-cache build-base
ARG GO111MODULE=on
ARG GOPROXY=https://goproxy.cn,direct
WORKDIR /app
COPY ./go.mod .
RUN ls
RUN  go mod download
