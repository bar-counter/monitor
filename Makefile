.PHONY: test check clean build dist all
#TOP_DIR := $(shell pwd)
# can change by env:ENV_CI_DIST_VERSION use and change by env:ENV_CI_DIST_MARK by CI
ENV_DIST_VERSION:=v2.1.0
ENV_DIST_MARK=

ROOT_NAME?=monitor

## MakeDocker.mk settings start
ROOT_OWNER ?=sinlov
ROOT_PARENT_SWITCH_TAG =1.23.4
# for image local build
INFO_TEST_BUILD_DOCKER_PARENT_IMAGE =golang
# for image running
INFO_BUILD_DOCKER_FROM_IMAGE =alpine:3.17
INFO_BUILD_DOCKER_FILE =Dockerfile
INFO_TEST_BUILD_DOCKER_FILE =build.dockerfile
## MakeDocker.mk settings end

## run info start
ENV_RUN_INFO_HELP_ARGS = -h
ENV_RUN_INFO_ARGS =
## run info end

## go test go-test.mk start
# ignore used not matching mode
# set ignore of test case like grep -v -E "vendor|go_fatal_error" to ignore vendor and go_fatal_error package
ENV_ROOT_TEST_INVERT_MATCH ?="vendor|go_fatal_error|robotn|shirou"
ifeq ($(OS),Windows_NT)
ENV_ROOT_TEST_LIST ?=./...
else
ENV_ROOT_TEST_LIST ?=$$(go list ./... | grep -v -E ${ENV_ROOT_TEST_INVERT_MATCH})
endif
# test max time
ENV_ROOT_TEST_MAX_TIME :=1m
## go test go-test.mk end

## clean args start
ENV_ROOT_BUILD_PATH =build
ENV_ROOT_LOG_PATH =logs/
## clean args end

## build args start
ENV_ROOT_BUILD_ENTRANCE =cmd/bar-counter-monitor/main.go
ENV_ROOT_BUILD_PATH =build
ENV_ROOT_BUILD_BIN_NAME =${ROOT_NAME}
ENV_ROOT_BUILD_BIN_PATH =${ENV_ROOT_BUILD_PATH}/${ENV_ROOT_BUILD_BIN_NAME}
## build args end

## build dist args start
# linux windows darwin  list as: go tool dist list
ENV_DIST_GO_OS =linux
# amd64 386
ENV_DIST_GO_ARCH =amd64
# mark for dist and tag helper
ENV_ROOT_MANIFEST_PKG_JSON? =package.json
ENV_ROOT_CHANGELOG_PATH ?=CHANGELOG.md
## build dist args end

include z-MakefileUtils/MakeBasicEnv.mk
include z-MakefileUtils/MakeDistTools.mk
include z-MakefileUtils/go-list.mk
include z-MakefileUtils/go-mod.mk
include z-MakefileUtils/go-test.mk
include z-MakefileUtils/go-test-integration.mk
include z-MakefileUtils/go-dist.mk
include z-MakefileUtils/MakeDocker.mk

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
	@echo "ENV_DIST_CODE_MARK                        ${ENV_DIST_CODE_MARK}"
	@echo "== project env info end =="

.PHONY: cleanBuild
cleanBuild:
	@$(RM) -r ${ENV_ROOT_BUILD_PATH}
	@echo "~> finish clean path: ${ENV_ROOT_BUILD_PATH}"

.PHONY: cleanLog
cleanLog:
	@$(RM) -r ${ENV_ROOT_LOG_PATH}
	@echo "~> finish clean path: ${ENV_ROOT_LOG_PATH}"

.PHONY: cleanTest
cleanTest: test.go.clean

.PHONY: cleanTestData
cleanTestData:
	$(info -> notes: remove folder [ testdata ] unable to match subdirectories)
	@$(RM) -r **/testdata
	@$(RM) -r **/**/testdata
	@$(RM) -r **/**/**/testdata
	@$(RM) -r **/**/**/**/testdata
	@$(RM) -r **/**/**/**/**/testdata
	@$(RM) -r **/**/**/**/**/**/testdata
	$(info -> finish clean folder [ testdata ])

.PHONY: clean
clean: cleanTest cleanBuild cleanLog
	@echo "~> clean finish"

.PHONY: cleanAll
cleanAll: clean
	@echo "~> clean all finish"

init:
	@echo "~> start init this project"
	@echo "-> check version"
	go version
	@echo "-> check env golang"
	go env
	@echo "~> you can use [ make help ] see more task"
	-go mod verify

.PHONY: dep
dep: go.mod.verify go.mod.download go.mod.tidy

.PHONY: style
style: go.mod.verify go.mod.tidy go.mod.fmt go.mod.lint.run

.PHONY: test
test: test.go

.PHONY: ci
ci: style go.mod.vet test

.PHONY: ci.test.benchmark
ci.test.benchmark: test.go.benchmark

.PHONY: ci.coverage.show
ci.coverage.show: test.go.coverage.show

.PHONY: ci.all
ci.all: ci ci.test.benchmark ci.coverage.show

.PHONY: buildMain
buildMain:
	@echo "-> start build local OS: ${PLATFORM} ${OS_BIT}"
ifeq ($(OS),Windows_NT)
	@go build -ldflags "-X main.buildID=${ENV_DIST_CODE_MARK}" -o ${ENV_ROOT_BUILD_BIN_PATH}.exe ${ENV_ROOT_BUILD_ENTRANCE}
	@echo "-> finish build out path: $(subst /,\,${ENV_ROOT_BUILD_BIN_PATH}).exe"
else
	@go build -ldflags "-X main.buildID=${ENV_DIST_CODE_MARK}" -o ${ENV_ROOT_BUILD_BIN_PATH} ${ENV_ROOT_BUILD_ENTRANCE}
	@echo "-> finish build out path: ${ENV_ROOT_BUILD_BIN_PATH}"
endif

.PHONY: devHelp
devHelp: export CI_DEBUG=false
devHelp: cleanBuild buildMain
ifeq ($(OS),Windows_NT)
	$(subst /,\,${ENV_ROOT_BUILD_BIN_PATH}).exe ${ENV_RUN_INFO_HELP_ARGS}
else
	${ENV_ROOT_BUILD_BIN_PATH} ${ENV_RUN_INFO_HELP_ARGS}
endif

.PHONY: dev
dev: export CI_DEBUG=true
dev: cleanBuild buildMain
ifeq ($(OS),Windows_NT)
	$(subst /,\,${ENV_ROOT_BUILD_BIN_PATH}).exe ${ENV_RUN_INFO_ARGS}
else
	${ENV_ROOT_BUILD_BIN_PATH} ${ENV_RUN_INFO_ARGS}
endif

.PHONY: runHelp
runHelp: export CI_DEBUG=false
runHelp:
	go run -v ${ENV_ROOT_BUILD_ENTRANCE} ${ENV_RUN_INFO_HELP_ARGS}

.PHONY: cloc
cloc:
	@echo "see: https://stackoverflow.com/questions/26152014/cloc-ignore-exclude-list-file-clocignore"
	cloc --exclude-list-file=.clocignore .

.PHONY: helpProjectRoot
helpProjectRoot:
	@echo "Help: Project root Makefile"
ifeq ($(OS),Windows_NT)
	@echo ""
	@echo "warning: other install make cli tools has bug, please use: scoop install main/make"
	@echo " run will at make tools version 4.+"
	@echo "windows use this kit must install tools blow:"
	@echo ""
	@echo "https://scoop.sh/#/apps?q=busybox&s=0&d=1&o=true"
	@echo "-> scoop install main/busybox"
	@echo "and"
	@echo "https://scoop.sh/#/apps?q=shasum&s=0&d=1&o=true"
	@echo "-> scoop install main/shasum"
	@echo ""
endif
	@echo "-- now build name: ${ROOT_NAME} version: ${ENV_DIST_VERSION}"
	@echo "-- distTestOS or distReleaseOS will out abi as: ${ENV_DIST_GO_OS} ${ENV_DIST_GO_ARCH} --"
	@echo ""
	@echo "~> make test                 - run test fast"
	@echo "~> make ci.all               - run CI tasks all"
	@echo "~> make ci.test.benchmark    - run CI tasks as test benchmark"
	@echo "~> make ci.coverage.show     - run CI tasks as test coverage and show"
	@echo "~> make ci                   - run CI tools tasks"
	@echo ""
	@echo "~> make env                  - print env of this project"
	@echo "~> make init                 - check base env of this project"
	@echo "~> make dep                  - check and install by go mod"
	@echo "~> make clean                - remove build binary file, log files, and testdata"
	@echo "~> make style                - run local code fmt and style check"
	@echo "~> make buildMain            - build binary file"
	@echo "~> make devHelp              - run local binary file args ${ENV_RUN_INFO_HELP_ARGS}"
	@echo "~> make dev                  - run local binary file args ${ENV_RUN_INFO_ARGS}"
	@echo "~> make runHelp              - run file with help args ${ENV_RUN_INFO_HELP_ARGS}"
	@echo ""

.PHONY: help
help: helpProjectRoot
	@echo "== show more help"
	@echo ""
	@echo "$$ make helpGoDist"
	@echo ""
	@echo "$$ make help.test.go.integration"
	@echo "$$ make help.test.go"
	@echo "$$ make help.go.list"
	@echo "$$ make help.go.mod"
	@echo ""
	@echo "-- more info see Makefile include --"

