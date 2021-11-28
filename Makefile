NOW = $(shell date -u '+%Y%m%d%I%M%S')

RELEASE_VERSION = v1.0.0

APP 			= app
SERVER_BIN  	= ./main/${APP}
RELEASE_ROOT 	= release
RELEASE_SERVER 	= release/${APP}
GIT_COUNT 		= $(shell git rev-list --all --count)
GIT_HASH        = $(shell git rev-parse --short HEAD)
RELEASE_TAG     = $(RELEASE_VERSION).$(GIT_COUNT).$(GIT_HASH)

all: start

.PHONY: build
build:
	@go build -ldflags "-w -s -X main.VERSION=$(RELEASE_TAG)" -o $(SERVER_BIN) ./main

.PHONY: start
start:
	@go run -ldflags "-X main.VERSION=$(RELEASE_TAG)" ./main/main.go web -c ./configs/config.toml -m ./configs/model.conf --menu ./configs/menu.yaml

.PHONY: swagger
swagger:
	@hash swag > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	@swag init --generalInfo ./main/main.go --output ./app/interfaces/api/swagger

.PHONY: wire
wire:
	@hash wire > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go install github.com/google/wire/cmd/wire@latest; \
	fi
	@wire gen ./injector

.PHONY: clean
clean:
	rm -rf data release $(SERVER_BIN) data

.PHONY: pack
pack: build
	rm -rf $(RELEASE_ROOT) && mkdir -p $(RELEASE_SERVER)
	cp -r $(SERVER_BIN) configs $(RELEASE_SERVER)
	cd $(RELEASE_ROOT) && tar -cvf $(APP).tar ${APP} && rm -rf ${APP}

.PHONY: lint
lint:
	@hash golangci-lint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1; \
	fi
	golangci-lint run ./...

.PHONY: docker-push
docker-push:
	time docker build --file ./Dockerfile --no-cache --tag seion/ddd-gin-admin .
	docker login
	docker push seion/ddd-gin-admin

.PHONY: skaffold-build
skaffold-build:
	skaffold build --file-output output.json

.PHONY: skaffold-dev
skaffold-dev:
	skaffold dev -v debug -p dev

.PHONY: tunnel-svc-with-minikube
tunnel-svc-with-minikube:
	minikube service ddd-gin-admin-web --url -n ddd-gin-admin
