APP_NAME=helium_web
PROJECT_NAME=helium
GIT_HASH ?= $(shell git rev-parse --short HEAD)
VERSION ?= 0.1
BUILD_OS ?= linux
BUILD_PATCH ?= develop
BUILD_VERSION ?= ${VERSION}.${BUILD_PATCH}

.ONESHELL:
.DEFAULT_GOAL:= build

.EXPORT_ALL_VARIABLES:
GO111MODULE=on
GOOS=${BUILD_OS}

.PHONY: install_lint lint test coverage

install_lint:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${HOME}/bin v1.21.0

lint:
	@echo "Running linters..."
	golangci-lint --version
	golangci-lint run

test:
	@echo "Running unit tests with coverage..."
	go test -v -cover -coverprofile=${APP_NAME}.coverprofile ./...

coverage: test
	go tool cover -html=${APP_NAME}.coverprofile

.PHONY: build run

build: GIT_TAG?=$(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match 2>/dev/null || git describe --all)
build: BUILD_TIME=$(shell date +%FT%T%z)
build: LDFLAGS=-w -s -X main.Name=${APP_NAME} -X main.Version=${BUILD_VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitTag=${GIT_TAG} -X main.GitHash=${GIT_HASH} -extldflags \"-static\"
build:
	@go build -v -ldflags "${LDFLAGS}"

run: build
run:
	./$(shell basename `pwd`)
