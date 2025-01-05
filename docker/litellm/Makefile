#
# make litellm
#

###
.DEFAULT_GOAL := help

.PHONY: help
help: Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

##
APP_VERSION ?= osi
APP_REPO_URL ?= https://github.com/qiangli/litellm.git
TARGET_DIR = local/litellm

##
clone:
	@if [ ! -d $(TARGET_DIR) ]; then \
		echo "Cloning repository..."; \
		git clone --recurse-submodules --branch $(APP_VERSION) $(APP_REPO_URL) $(TARGET_DIR); \
	else \
		echo "Directory exists."; \
	fi; \

checkout:
	@cd $(TARGET_DIR) && git stash && git pull origin $(APP_VERSION)

.PHONY: clone checkout

##
build: clone checkout ## Build LiteLLM
	@docker compose build

up: net ## Start LiteLLM
	@docker compose up -d

down:  ## Stop LiteLLM
	@docker compose down

net: ## Create openland network
	@docker network inspect openland >/dev/null 2>&1 \
		|| docker network create openland

.PHONY: build up down net
###
