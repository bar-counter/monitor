.PHONY: test check clean build dist all
#TOP_DIR := $(shell pwd)
# can change by env:ENV_CI_DIST_VERSION use and change by env:ENV_CI_DIST_MARK by CI
ENV_DIST_VERSION:=v2.1.0
ENV_DIST_MARK=

ROOT_NAME?=monitor

## run info start
ENV_RUN_INFO_HELP_ARGS= -h
ENV_RUN_INFO_ARGS=
## run info end

## build dist env start
# change to other build entrance
ENV_ROOT_BUILD_ENTRANCE = main.go
ENV_ROOT_BUILD_BIN_NAME = $(ROOT_NAME)
ENV_ROOT_BUILD_PATH = build
ENV_ROOT_BUILD_BIN_PATH = $(ENV_ROOT_BUILD_PATH)/$(ENV_ROOT_BUILD_BIN_NAME)
ENV_ROOT_LOG_PATH = log/
# linux windows darwin  list as: go tool dist list
ENV_DIST_GO_OS=linux
# amd64 386
ENV_DIST_GO_ARCH=amd64
# mark for dist and tag helper
ENV_ROOT_MANIFEST_PKG_JSON?=package.json
ENV_ROOT_MAKE_FILE?=Makefile
ENV_ROOT_CHANGELOG_PATH?=CHANGELOG.md
## build dist env end

## go test MakeGoTest.mk start
# ignore used not matching mode
# set ignore of test case like grep -v -E "vendor|go_fatal_error" to ignore vendor and go_fatal_error package
ENV_ROOT_TEST_INVERT_MATCH ?= "vendor|go_fatal_error|robotn|shirou|go_robot"
ifeq ($(OS),Windows_NT)
ENV_ROOT_TEST_LIST?=./...
else
ENV_ROOT_TEST_LIST?=$$(go list ./... | grep -v -E ${ENV_ROOT_TEST_INVERT_MATCH})
endif
# test max time
ENV_ROOT_TEST_MAX_TIME:=1m
## go test MakeGoTest.mk en

include z-MakefileUtils/MakeBasicEnv.mk
include z-MakefileUtils/MakeDistTools.mk
include z-MakefileUtils/MakeGoMod.mk
include z-MakefileUtils/MakeGoTest.mk
include z-MakefileUtils/MakeGoDist.mk

all: env

env: distEnv
	@echo "== project env info start =="
	@echo ""
	@echo "test info"
	@echo "ENV_ROOT_TEST_LIST                        ${ENV_ROOT_TEST_LIST}"
	@echo ""
	@echo "ROOT_NAME                                 ${ROOT_NAME}"
	@echo "ENV_DIST_VERSION                          ${ENV_DIST_VERSION}"
	@echo "ENV_ROOT_CHANGELOG_PATH                   ${ENV_ROOT_CHANGELOG_PATH}"
	@echo ""
	@echo "ENV_ROOT_BUILD_ENTRANCE                   ${ENV_ROOT_BUILD_ENTRANCE}"
	@echo "ENV_ROOT_BUILD_PATH                       ${ENV_ROOT_BUILD_PATH}"
ifeq ($(OS),Windows_NT)
	@echo "ENV_ROOT_BUILD_BIN_PATH                   $(subst /,\,${ENV_ROOT_BUILD_BIN_PATH}).exe"
else
	@echo "ENV_ROOT_BUILD_BIN_PATH                   ${ENV_ROOT_BUILD_BIN_PATH}"
endif
	@echo "ENV_DIST_GO_OS                            ${ENV_DIST_GO_OS}"
	@echo "ENV_DIST_GO_ARCH                          ${ENV_DIST_GO_ARCH}"
	@echo ""
	@echo "ENV_DIST_MARK                             ${ENV_DIST_MARK}"
	@echo "== project env info end =="

cleanBuild:
	-@$(RM) -r ${ENV_ROOT_BUILD_PATH}
	@echo "~> finish clean path: ${ENV_ROOT_BUILD_PATH}"

cleanLog:
	-@$(RM) -r ${ENV_ROOT_LOG_PATH}
	@echo "~> finish clean path: ${ENV_ROOT_LOG_PATH}"

cleanTestData:
	$(info -> notes: remove folder [ testdata ] unable to match subdirectories)
	@$(RM) coverage.txt
	@$(RM) -r **/testdata
	@$(RM) -r **/**/testdata
	@$(RM) -r **/**/**/testdata
	@$(RM) -r **/**/**/**/testdata
	@$(RM) -r **/**/**/**/**/testdata
	@$(RM) -r **/**/**/**/**/**/testdata
	$(info -> finish clean folder [ testdata ])

clean: cleanTestData cleanBuild cleanLog
	@echo "~> clean finish"

cleanAll: clean cleanAllDist
	@echo "~> clean all finish"

init:
	@echo "~> start init this project"
	@echo "-> check version"
	go version
	@echo "-> check env golang"
	go env
	@echo "~> you can use [ make help ] see more task"
	-go mod verify

dep: modVerify modDownload modTidy modVendor
	@echo "-> just check depends below"

ci: modTidy modVerify modFmt modVet modLintRun test

exampleStatus: dep
	@echo "=> run dev example/status/statusdemo"
	@go run example/status/statusdemo.go

exampleDebug: dep
	@echo "=> run dev example/debug/debugdemo"
	@go run example/debug/debugdemo.go

examplePprof: dep
	@echo "=> run dev example/pprof/pprofdemo"
	@go run example/pprof/pprofdemo.go

helpProjectRoot:
	@echo "Help: Project root Makefile"
	@echo "-- now build name: ${ROOT_NAME} version: ${ENV_DIST_VERSION}"
	@echo "-- distTestOS or distReleaseOS will out abi as: ${ENV_DIST_GO_OS} ${ENV_DIST_GO_ARCH} --"
	@echo ""
	@echo "~> make env                 - print env of this project"
	@echo "~> make init                - check base env of this project"
	@echo "~> make dep                 - check and install by go mod"
	@echo "~> make clean               - remove build binary file, log files, and testdata"
	@echo "~> make test                - run test case ignore --invert-match by config"
	@echo "~> make testCoverage        - run test coverage case ignore --invert-match by config"
	@echo "~> make testCoverageBrowser - see coverage at browser --invert-match by config"
	@echo "~> make testBenchmark       - run go test benchmark case all"
	@echo "~> make ci                  - run CI tools tasks"

help: helpGoMod helperGoTest helpDist helpProjectRoot
	@echo ""
	@echo "-- more info see Makefile include: MakeGoMod.mk MakeGoTest.mk MakeGoDist.mk --"