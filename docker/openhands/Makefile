#
# make openhands
#

###
.DEFAULT_GOAL := help

.PHONY: help
help: Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

##
APP_VERSION ?= main
APP_REPO_URL ?= https://github.com/qiangli/OpenHands.git
# APP_REPO_URL ?= https://github.com/All-Hands-AI/OpenHands.git
TARGET_DIR = local/openhands

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
build: clone ## Build mage
	@docker compose build

up: net ## Start
	@docker compose up -d

down:  ## Stop
	@docker compose down

net: ## Create network
	@docker network inspect openland >/dev/null 2>&1 \
		|| docker network create openland

cli: ## Build openhands cli/wrapper for linux/macos/windows
	@echo "Building openhands cli with docker..."
	@cd cli && ./build.sh

clean: ## Clean build artifacts
	@rm -rf cli/dist/*

.PHONY: clone build up down net
.PHONY: cli clean
###
