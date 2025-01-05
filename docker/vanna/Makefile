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
	else \
		echo "Directory exists."; \
	fi;

checkout:
	@cd $(TARGET_DIR) && git stash && git checkout $(APP_VERSION) && git apply $(PATCH)

# debug
setup-vscode:
	@cp .vscode/launch.example.json .vscode/launch.json && \
	touch .env && \
	echo "Please update .env file with database credentials"

.PHONY: clone checkout setup-vscode

##
build: clone checkout ## Build docker images
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

setup: ## setup venv and install requirements
	@setup.sh

.PHONY: build up down net proxy setup
##