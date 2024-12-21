#
# make aider
#

###
.DEFAULT_GOAL := help

.PHONY: help
help: Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

##
APP_VERSION ?= main
APP_REPO_URL ?= https://github.com/Aider-AI/aider.git
TARGET_DIR = local/aider

##
clone:
	@if [ ! -d $(TARGET_DIR) ]; then \
		echo "Cloning repository..."; \
		git clone --branch $(APP_VERSION) $(APP_REPO_URL) $(TARGET_DIR); \
	else \
		echo "Directory exists."; \
	fi;

checkout:
	@cd $(TARGET_DIR) && \
	git checkout $(APP_VERSION) && \
	git stash && \
	git pull origin $(APP_VERSION)

.PHONY: clone checkout

##
build: clone checkout ## Build aider
	@docker compose --profile base build aider-base
	@docker compose build

up: net ## Start aider
	@docker compose up -d

down:  ## Stop aider
	@docker compose down

net: ## Create openland network
	@docker network inspect openland >/dev/null 2>&1 \
		|| docker network create openland

.PHONY: build up down net

#
app: ## setup aider (python)
	@build-app.sh
.PHONY: app
###
