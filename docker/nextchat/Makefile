###

#
# make nextchat
#
.DEFAULT_GOAL := help

.PHONY: help
help: Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

##
APP_VERSION ?= main
APP_REPO_URL ?= https://github.com/ChatGPTNextWeb/ChatGPT-Next-Web.git
TARGET_DIR = local/nextchat

##
clone:
	@if [ ! -d $(TARGET_DIR) ]; then \
		echo "Cloning repository..."; \
		git clone --branch $(APP_VERSION) $(APP_REPO_URL) $(TARGET_DIR); \
	else \
		echo "Directory exists."; \
	fi; \

checkout:
	@cd $(TARGET_DIR) && git checkout $(APP_VERSION)

##
build: clone checkout ## Build nextchat image
	@docker buildx bake

up: ## Start nextchat
	@docker compose up -d

down: ## Stop nextchat
	@docker compose down

net: ## Create network
	@docker network inspect openland >/dev/null 2>&1 \
		|| docker network create openland

.PHONY: clone build up down net
##