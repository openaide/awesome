#
# make danswer
#

.DEFAULT_GOAL := help

.PHONY: help
help: Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

##
APP_VERSION ?= main
APP_REPO_URL ?= https://github.com/danswer-ai/danswer.git
TARGET_DIR = local/danswer

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
# https://docs.danswer.dev/quickstart
# docker compose -f docker-compose.dev.yml -p awesome up -d --build --force-recreate
# docker compose -f docker-compose.dev.yml -p awesome down
build: clone ## Build mage
	@cd local/danswer/deployment/docker_compose && \
		docker compose -f docker-compose.dev.yml build api_server web_server inference_model_server

up: net ## Start
	@docker compose up -d

down: ## Stop
	@docker compose down

net: ## Create openland network
	@docker network inspect openland >/dev/null 2>&1 \
		|| docker network create openland

.PHONY: clone build up down net
###