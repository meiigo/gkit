PWD=$(shell pwd)
GKIT=$(GOPATH)/pkg/mod/github.com/meiigo/gkit
API_PROTO_FILES=$(shell find api -name *.proto)

## 将 phony 对象当做命令执行而非执行文件执行
.PHONY: info allure

APP_NAME=blog
APP_VERSION=1.0.0

## 项目通用变量
# git 参数
GIT_VERSION=$(shell git tag | grep $(APP_NAME)-v | sort -r --version-sort | head -n1)
GIT_COMMIT_ID=$(shell git rev-parse --short=8 HEAD)
GIT_BRANCH=$(shell git branch | sed -n -e 's/^\* \(.*\)/\1/p')
GIT_MASTER_DIFF=$(shell git log --no-merges -q origin/master ^$(GIT_BRANCH) | grep commit | awk '{print $$2}')
GIT_MASTER_DIFF_COUNT=$(shell echo $(GIT_MASTER_DIFF) | wc -w | sed 's/^[[:space:]]*//g')

# go 环境变量
GO111MODULE=on
GOPROXY=https://goproxy.cn,direct
GOARCH=amd64
GOOS=linux
CGO_ENABLED=0

.PHONY: proto
# generate grpc code
proto:
	protoc --proto_path=. \
           --proto_path=$(GKIT)/third_party \
           --go_out=paths=source_relative:. \
           --go-grpc_out=paths=source_relative:. \
           $(API_PROTO_FILES)

.PHONY: build
# build
build:
	go build -o bin/$(APP_NAME) -a -ldflags '-w -s' -tags=jsoniter cmd/*.go

.PHONY: run
run:
	go run cmd/*.go -conf $(PWD)/config/local

test:
	@echo $(API_PROTO_FILES)

# test url
# http://127.0.0.1:8080/_ping
# http://127.0.0.1:16060/monitor/env
# http://127.0.0.1:16060/monitor/metrics
# http://127.0.0.1:16060/monitor/debug/pprof/
