#
# make screenshot-to-code
#

###
.DEFAULT_GOAL := help

.PHONY: help
help: Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

##
APP_VERSION ?= main
APP_REPO_URL ?= https://github.com/abi/screenshot-to-code.git
TARGET_DIR = local/screenshot-to-code

PATCH=$(PWD)/src.patch

##
clone:
	@if [ ! -d $(TARGET_DIR) ]; then \
		echo "Cloning repository..."; \
		git clone --branch $(APP_VERSION) $(APP_REPO_URL) $(TARGET_DIR); \
	else \
		echo "Directory exists."; \
	fi; \

checkout:
	@cd $(TARGET_DIR) && \
		git checkout $(APP_VERSION) && \
		git pull origin $(APP_VERSION)

patch:
	@cd $(TARGET_DIR) && git reset --hard && git apply $(PATCH)

# setup debug
setup:
	@mkdir -p $(TARGET_DIR)/backend/.vscode && \
	cp .vscode/launch.example.json $(TARGET_DIR)/backend/.vscode/launch.json && \
	mkdir -p $(TARGET_DIR)/frontend/.vscode && \
	cp .vscode/launch.example.json $(TARGET_DIR)/frontend/.vscode/launch.json


.PHONY: clone checkout patch setup

##
build: clone checkout patch ## Build
	@docker compose build

rebuild:
	@docker compose build --no-cache

up: net ## Start
	@docker compose up -d

down: ## Stop
	@docker compose down

net: ## Create openland network
	@docker network inspect openland >/dev/null 2>&1 \
		|| docker network create openland

.PHONY: build rebuild up down net
###
