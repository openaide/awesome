#
# make anythingllm
#

.DEFAULT_GOAL := help

.PHONY: help
help: Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

##
# APP_VERSION ?= master
APP_VERSION ?= v3.1.4
APP_REPO_URL ?= https://github.com/assafelovic/gpt-researcher.git
TARGET_DIR = local/gpt-researcher

PATCH=$(PWD)/host.patch

##
clone:
	@if [ ! -d $(TARGET_DIR) ]; then \
		echo "Cloning repository..."; \
		git clone --branch $(APP_VERSION) $(APP_REPO_URL) $(TARGET_DIR); \
		cd $(TARGET_DIR) && git checkout -b $(APP_VERSION) && git apply $(PATCH); \
	else \
		echo "Directory exists."; \
	fi;

##
build: clone build-server build-nextjs ## Build mage
	
build-server: ## Build server
	@docker buildx bake gptr-server

build-nextjs: ## Build nextjs
	@docker buildx bake gptr-nextjs

up: net ## Start
	@docker compose up -d

down: ## Stop
	@docker compose down

net: ## Create openland network
	@docker network inspect openland >/dev/null 2>&1 \
		|| docker network create openland

.PHONY: clone build build-server build-nextjs up down net
###