GO_EXECUTABLE ?= go
BUILD_VERSION ?= $(shell git describe --tags)
BUILD_TIME = `date +%FT%T%z`
BUILD_NAME = go-cd
MAIN_FILE = go-cd.go
LIST_OF_FILES = $(shell ${GO_EXECUTABLE} list ./... | grep -v /vendor/ | grep -v /src/)

export PATH := $(PATH):$(GOPATH)/bin

init:
	${GO_EXECUTABLE} get github.com/Masterminds/glide
	${GOPATH}/bin/glide ${GLIDEOPS} install


build:
	${GO_EXECUTABLE} build \
	-o build/${BUILD_NAME} \
	-ldflags="-X main.Version=${BUILD_VERSION} -X main.BuildTime=${BUILD_TIME}" \
	.

run: build
	./build/${BUILD_NAME}

test:
	${GO_EXECUTABLE} vet -tests ${LIST_OF_FILES}
	${GO_EXECUTABLE} test -gcflags=-l -race -cover -bench . ${LIST_OF_FILES}

full-test:
	${GO_EXECUTABLE} get github.com/jgautheron/goconst
	${GO_EXECUTABLE} get github.com/alecthomas/gocyclo
	${GO_EXECUTABLE} get mvdan.cc/interfacer
	${GO_EXECUTABLE} get github.com/walle/lll/cmd/lll
	${GO_EXECUTABLE} get github.com/mdempsky/unconvert
	${GO_EXECUTABLE} get mvdan.cc/unparam
	${GO_EXECUTABLE} get honnef.co/go/tools/cmd/unused
	unused ${LIST_OF_FILES}
	unparam -tests=false ${LIST_OF_FILES}
	unconvert ${LIST_OF_FILES}
	lll -g -l 140
	interfacer ${LIST_OF_FILES}
	gofmt -l -s .
	gocyclo -over 15 -avg .
	goconst -min-occurrences 3 -min-length 3 -ignore-tests .

.PHONY: init build test full-test
