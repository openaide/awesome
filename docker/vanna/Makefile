#
# make vanna
#
.DEFAULT_GOAL := help

.PHONY: help
help: Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

##
TARGET_DIR = local/vanna

APP_REPO_URL ?= https://github.com/vanna-ai/vanna.git
VANNA_STREAMLIT = https://github.com/vanna-ai/vanna-streamlit.git
VANNA_FLASK = https://github.com/vanna-ai/vanna-flask.git
VANNA_NEXTJS_FLASK = https://github.com/vanna-ai/vanna-nextjs-flask.git

##
clone:
	@if [ ! -d $(TARGET_DIR) ]; then \
		echo "Cloning repository..."; \
		git clone --branch main $(APP_REPO_URL) $(TARGET_DIR)/vanna; \
		git clone --branch main $(VANNA_STREAMLIT) $(TARGET_DIR)/vanna-streamlit; \
		git clone --branch main $(VANNA_FLASK) $(TARGET_DIR)/vanna-flask; \ 
		git clone --branch main $(VANNA_NEXTJS_FLASK) $(TARGET_DIR)/vanna-nextjs-flask; \
	else \
		echo "Directory exists."; \
	fi; \

build: clone ## Build docker images
	@docker buildx bake

up: ## Start
	@docker compose up -d

down: ## Stop
	@docker compose down

net: ## Create openland network
	@docker network inspect openland >/dev/null 2>&1 \
		|| docker network create openland

.PHONY: clone build up down net
##