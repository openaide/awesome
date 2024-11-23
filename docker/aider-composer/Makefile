#
# make aider-composer
#

###
.DEFAULT_GOAL := help

.PHONY: help
help: Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

##
# APP_VERSION ?= v1.3.0
APP_VERSION ?= master
APP_REPO_URL ?= https://github.com/lee88688/aider-composer.git
TARGET_DIR = local/aider-composer

OS := $(shell uname)
ARCH := $(shell uname -m)

VSCE_TARGET := $(if $(findstring Darwin,$(OS)),darwin,linux)-$(if $(findstring arm64,$(ARCH)),arm64,x64)

info:
	@echo "VSCE_TARGET: $(VSCE_TARGET)"
	@echo "APP_REPO_URL: $(APP_REPO_URL)"
	@echo "APP_VERSION: $(APP_VERSION)"
	@echo "TARGET_DIR: $(TARGET_DIR)"

.PHONY: info

##
clone:
	@if [ ! -d $(TARGET_DIR) ]; then \
		echo "Cloning repository..."; \
		git clone --branch $(APP_VERSION) $(APP_REPO_URL) $(TARGET_DIR); \
	else \
		echo "Directory exists."; \
	fi; \
	cd $(TARGET_DIR) && git pull origin $(APP_VERSION)

##
build: vsce

build-service: clone ## Build aider-composer chat service image
	@docker buildx bake --progress plain

up: net ## Start aider-composer chat service
	@docker compose up -d

down:  ## Stop aider-composer chat service
	@docker compose down

net: ## Create openland network
	@docker network inspect openland >/dev/null 2>&1 \
		|| docker network create openland

.PHONY: clone build-service up down net

##
vsce: clone clean ## Build aider-composer vscode extension
	@docker buildx build --build-arg VSCE_TARGET=$(VSCE_TARGET) --progress plain --target vsce --output type=local,dest=./local/extension --file ./Dockerfile $(TARGET_DIR)

clean:
	@rm -rf ./local/extension

.PHONY: build vsce clean
###
