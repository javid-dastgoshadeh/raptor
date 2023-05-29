NAME = Oauth service

BINARY = service
ROOTDIR ?= ${PWD}
VERSION := `git describe --tags`
BINDIR ?= "${ROOTDIR}/bin"
EXT = ""
OUTPUT = "${BINDIR}/${BINARY}-${VERSION}${EXT}"
INSTALL_PATH ?= ${GOBIN}
BUILD=`date +%FT%T%z`
LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION} -X main.Build=${BUILD}"
OS ?= `uname -s | tr '[:upper:]' '[:lower:]'`
ARCH ?= amd64

export GOBIN = "${ROOTDIR}/bin"

ifeq (${OS}, windows)
	EXT := .exe
endif

.PHONY: build install clean

default: build

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

list:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m\n", $$1}'

version: ## Get app's current version
	@echo $(VERSION)

dep: ## Install dependencies
	dep ensure

build: ## Build all targets in $ROOTDIR into $BINDIR
	GOOS=${OS} GOARCH=${ARCH} go build ${LDFLAGS} -o ${OUTPUT}

deploy: build
deploy: ## Deploy app in remote server
	scp ${OUTPUT} javiddastgoshadeh.com:/opt/javid/services/identity

test: ## Test app
	go test

windows: OS := windows
windows: EXT := ".exe"
windows: build ## Build for windows

linux: OS := linux
linux: build ## Build for linux

darwin: OS := darwin
darwin: build ## Build for darwin

rebuild: clean build ## Clean last build and build again

install: build ## Install app in $INSTALL_PATH
	mkdir -p "$(INSTALL_PATH)"
	cp -f $(OUTPUT) "${INSTALL_PATH}/${BINARY}${EXT}"
	cp -f "${ROOTDIR}/config.json" $(INSTALL_PATH)

clean: ## Remove last build output
	rm -rf $(OUTPUT)

uninstall: ## Remove all files created by install target
	rm -f "${INSTALL_PATH}/${BINARY}${EXT}"
	rm -rf "${INSTALL_PATH}/config.json"
