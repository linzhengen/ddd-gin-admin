.PHONY: start build

NOW = $(shell date -u '+%Y%m%d%I%M%S')

RELEASE_VERSION = v1.0.0

APP 			= main
SERVER_BIN  	= ./main/${APP}
RELEASE_ROOT 	= release
RELEASE_SERVER 	= release/${APP}
GIT_COUNT 		= $(shell git rev-list --all --count)
GIT_HASH        = $(shell git rev-parse --short HEAD)
RELEASE_TAG     = $(RELEASE_VERSION).$(GIT_COUNT).$(GIT_HASH)

all: start

build:
	@go build -ldflags "-w -s -X main.VERSION=$(RELEASE_TAG)" -o $(SERVER_BIN) ./main

start:
	@go run -ldflags "-X main.VERSION=$(RELEASE_TAG)" ./main/main.go web -c ./configs/config.toml -m ./configs/model.conf --menu ./configs/menu.yaml

swagger:
	@swag init --generalInfo ./main/main.go --output ./interfaces/swagger

wire:
	@wire gen ./injector

clean:
	rm -rf data release $(SERVER_BIN) data

pack: build
	rm -rf $(RELEASE_ROOT) && mkdir -p $(RELEASE_SERVER)
	cp -r $(SERVER_BIN) configs $(RELEASE_SERVER)
	cd $(RELEASE_ROOT) && tar -cvf $(APP).tar ${APP} && rm -rf ${APP}
