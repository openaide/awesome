#
# make smolagents
#
.DEFAULT_GOAL := help

.PHONY: help
help: Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

##
# APP_VERSION ?= main
APP_VERSION ?= v1.0.0
APP_REPO_URL ?= https://github.com/huggingface/smolagents.git
TARGET_DIR = local/smolagents

##
clone:
	@if [ ! -d $(TARGET_DIR) ]; then \
		echo "Cloning repository..."; \
		git clone $(APP_REPO_URL) $(TARGET_DIR); \
	else \
		echo "Directory exists."; \
	fi;

checkout:
	@cd $(TARGET_DIR) && git checkout -b $(APP_VERSION) $(APP_VERSION)

.PHONY: clone checkout

fix: 
	@cp requirements.txt $(TARGET_DIR)

##
build: clone fix ## Build docker images
	@docker buildx bake

up: ## Start
	@docker compose up -d

down: ## Stop
	@docker compose down

tail:
	@docker compose -f compose.yml -f compose.override.yml logs -f

net: ## Create openland network
	@docker network inspect openland >/dev/null 2>&1 \
		|| docker network create openland

setup: ## Setup
	@setup.sh

.PHONY: build up down net setup
##