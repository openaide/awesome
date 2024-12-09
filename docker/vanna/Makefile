#
# make vanna
#
.DEFAULT_GOAL := help

.PHONY: help
help: Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

##
APP_VERSION ?= main
APP_REPO_URL ?= https://github.com/vanna-ai/vanna.git
TARGET_DIR = local/vanna

PATCH=$(PWD)/assets.patch

##
clone:
	@if [ ! -d $(TARGET_DIR) ]; then \
		echo "Cloning repository..."; \
		git clone --branch $(APP_VERSION) $(APP_REPO_URL) $(TARGET_DIR); \
		cd $(TARGET_DIR) && git checkout $(APP_VERSION) && git apply $(PATCH); \
	else \
		echo "Directory exists."; \
	fi;

build: clone ## Build docker images
	@docker buildx bake

up: ## Start
	@docker compose up -d

down: ## Stop
	@docker compose down

net: ## Create openland network
	@docker network inspect openland >/dev/null 2>&1 \
		|| docker network create openland

proxy: ## build proxy (go)
	@build-proxy.sh

.PHONY: clone build up down net proxy
##