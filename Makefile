.PHONY: all get build clean run
.DEFAULT_GOAL: $(BIN_FILE)

PROJECT_NAME = kratgo

BIN_DIR = ./bin
BIN_FILE = $(PROJECT_NAME)

INTERNAL_DIR = ./internal

PROXY_DIR = $(INTERNAL_DIR)/proxy
CMD_DIR = ./cmd/$(PROJECT_NAME)
CONFIG_DIR = ./config/

# Get version constant
VERSION := $(shell cat $(PROXY_DIR)/const.go | grep "const version = " | awk '{print $$NF}' | sed -e 's/^.//' -e 's/.$$//')
BUILD := $(shell git rev-parse HEAD)

# Use linker flags to provide version/build settings to the binary
LDFLAGS=-ldflags "-s -w -X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"


default: get build

check-path:
ifndef GOPATH
	@echo "FATAL: you must declare GOPATH environment variable, for more"
	@echo "       details, please check"
	@echo "       http://golang.org/doc/code.html#GOPATH"
	@exit 1
endif

get: check-path
	@echo "[*] Downloading dependencies..."
	go get
	@echo "[*] Finish..."

vendor: check-path
	@go mod vendor

build:
	@echo "[*] Building $(PROJECT_NAME)..."
	go build $(LDFLAGS) -o $(BIN_DIR)/$(BIN_FILE) $(CMD_DIR)/...
	@echo "[*] Finish..."

test:
	go test -v -race -cover ./...

bench:
	go test -bench=. -benchmem ./...

run: build
	$(BIN_DIR)/$(BIN_FILE) -config ./config/kratgo-dev.conf.yml

install:
	mkdir -p /etc/kratgo/
	cp $(BIN_DIR)/$(BIN_FILE) /usr/local/bin/
	cp $(CONFIG_DIR)/kratgo.conf.yml /etc/kratgo/

uninstall:
	rm -rf /usr/local/bin/$(BIN_FILE)
	rm -rf /etc/kratgo/

clean:
	rm -rf bin/
	rm -rf vendor/

docker_build:
	docker build -t kratgo .
