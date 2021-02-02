.PHONY: test check clean build dist all

TOP_DIR := $(shell pwd)

# ifeq ($(FILE), $(wildcard $(FILE)))
# 	@ echo target file not found
# endif

DIST_VERSION := 1.1.0
DIST_OS := linux
DIST_ARCH := amd64

ROOT_BUILD_PATH ?= ./build
ROOT_DIST ?= ./dist
ROOT_TEST_DIST_PATH ?= $(ROOT_DIST)/test/$(DIST_VERSION)
ROOT_TEST_OS_DIST_PATH ?= $(ROOT_DIST)/$(DIST_OS)/test/$(DIST_VERSION)
ROOT_REPO_DIST_PATH ?= $(ROOT_DIST)/release/$(DIST_VERSION)
ROOT_REPO_OS_DIST_PATH ?= $(ROOT_DIST)/$(DIST_OS)/release/$(DIST_VERSION)

ROOT_LOG_PATH ?= ./log

# can use as https://goproxy.io/ https://gocenter.io https://mirrors.aliyun.com/goproxy/
ENV_GO_PROXY ?= https://goproxy.cn/

# include
include MakeGoMod.mk
include MakeGoTravis.mk

checkEnvGo:
ifndef GOPATH
	@echo Environment variable GOPATH is not set
	exit 1
endif

# check must run environment
init:
	@echo "~> start init this project"
	@echo "-> check version"
	go version
	@echo "-> check env golang"
	go env
	-GOPROXY="$(ENV_GO_PROXY)" GO111MODULE=on go mod download
	-GOPROXY="$(ENV_GO_PROXY)" GO111MODULE=on go mod vendor
	@echo "~> you can use [ make help ] see more task"

cleanBuild:
	@if [ -d ${ROOT_BUILD_PATH} ]; then rm -rf ${ROOT_BUILD_PATH} && echo "~> cleaned ${ROOT_BUILD_PATH}"; else echo "~> has cleaned ${ROOT_BUILD_PATH}"; fi

cleanLog:
	@if [ -d ${ROOT_LOG_PATH} ]; then rm -rf ${ROOT_LOG_PATH} && echo "~> cleaned ${ROOT_LOG_PATH}"; else echo "~> has cleaned ${ROOT_LOG_PATH}"; fi

clean: cleanBuild cleanLog
	@echo "~> clean finish"

checkTestDistPath:
	@if [ ! -d ${ROOT_TEST_DIST_PATH} ]; then mkdir -p ${ROOT_TEST_DIST_PATH} && echo "~> mkdir ${ROOT_TEST_DIST_PATH}"; fi

checkTestOSDistPath:
	@if [ ! -d ${ROOT_TEST_OS_DIST_PATH} ]; then mkdir -p ${ROOT_TEST_OS_DIST_PATH} && echo "~> mkdir ${ROOT_TEST_OS_DIST_PATH}"; fi

checkReleaseDistPath:
	@if [ ! -d ${ROOT_REPO_DIST_PATH} ]; then mkdir -p ${ROOT_REPO_DIST_PATH} && echo "~> mkdir ${ROOT_REPO_DIST_PATH}"; fi

checkReleaseOSDistPath:
	@if [ ! -d ${ROOT_REPO_OS_DIST_PATH} ]; then mkdir -p ${ROOT_REPO_OS_DIST_PATH} && echo "~> mkdir ${ROOT_REPO_OS_DIST_PATH}"; fi

buildMain: dep cleanBuild
	@go build -o $(ROOT_BUILD_PATH)/main example/pprof/pprofdemo.go

buildARCH: dep cleanBuild
	@GOOS=$(DIST_OS) GOARCH=$(DIST_ARCH) go build -o $(ROOT_BUILD_PATH)/main main.go

dev: buildMain
	-$(ROOT_BUILD_PATH)/main

helpProjectRoot:
	@echo "Help: Project root Makefile"
	@echo "~> make init  - check base env of this project"
	@echo "~> make clean - remove binary file and log files"
	@echo "~> make dev   - run example example/pprof/pprofdemo.go"

help: helpGoTravis helpGoMod helpProjectRoot
	@echo ""
	@echo "-- more info see Makefile include: MakeGoMod.mk --"
