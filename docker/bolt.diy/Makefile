#
# make bolt.diy
#

###
.DEFAULT_GOAL := help

.PHONY: help
help: Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

##
APP_VERSION ?= main
APP_REPO_URL ?= https://github.com/stackblitz-labs/bolt.diy.git
TARGET_DIR = local/bolt.diy

PATCH=$(PWD)/models.patch

##
clone:
	@if [ ! -d $(TARGET_DIR) ]; then \
		echo "Cloning repository..."; \
		git clone --branch $(APP_VERSION) $(APP_REPO_URL) $(TARGET_DIR); \
	else \
		echo "Directory exists."; \
	fi; \
	cd $(TARGET_DIR) && git pull origin $(APP_VERSION)

checkout: clone
	@cd $(TARGET_DIR) && git checkout $(APP_VERSION)

patch: checkout
	@cd $(TARGET_DIR) && git reset --hard && git apply $(PATCH)

.PHONY: clone checkout patch

##
build: patch ## Build
	@docker build ./local/bolt.diy --target bolt-ai-development

rebuild:
	@docker build ./local/bolt.diy --target bolt-ai-development --no-cache	

up: net ## Start
	@docker compose up -d

down: ## Stop
	@docker compose down

net: ## Create openland network
	@docker network inspect openland >/dev/null 2>&1 \
		|| docker network create openland

.PHONY: build rebuild up down net
###
