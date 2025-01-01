#
# make gpt-researcher
#

.DEFAULT_GOAL := help

.PHONY: help
help: Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

##
# APP_VERSION ?= master
# APP_VERSION ?= v3.1.4
APP_VERSION ?= v.3.1.7
APP_REPO_URL ?= https://github.com/assafelovic/gpt-researcher.git
TARGET_DIR = local/gpt-researcher

PATCH=$(PWD)/src.patch

##
clone:
	@if [ ! -d $(TARGET_DIR) ]; then \
		echo "Cloning repository..."; \
		git clone --branch $(APP_VERSION) $(APP_REPO_URL) $(TARGET_DIR); \
		cd $(TARGET_DIR) && git checkout -b $(APP_VERSION); \
	else \
		echo "Directory exists."; \
	fi;

checkout:
	@cd $(TARGET_DIR) && \
	git fetch origin && \
	git stash && \
	git checkout $(APP_VERSION) && \
	git apply $(PATCH)

.PHONY: clone checkout

##
build: clone checkout build-server build-nextjs ## Build mage
	
build-server: ## Build server
	@docker buildx bake gptr-server

build-nextjs: ## Build nextjs
	@docker buildx bake gptr-nextjs

up: net ## Start
	@docker compose up -d

down: ## Stop
	@docker compose down

net: ## Create openland network
	@docker network inspect openland >/dev/null 2>&1 \
		|| docker network create openland

.PHONY: build build-server build-nextjs up down net

setup: ## setup (python)
	@setup.sh

searxng: start-searxng ## Start standalone SearxNG API Retriever
start-searxng:
	@docker compose -f compose.searxng.yml up -d

stop-searxng: ## Stop standalone SearxNG API Retriever
	@docker compose -f compose.searxng.yml down

test: searxng ## Run llm/retriever tests
	@test.sh

cli: searxng ## Run CLI
	@cli.sh

app: searxng ## Run custom app
	@app.sh

tail-searxng: ## Tail logs
	@docker compose -f compose.searxng.yml logs -f

.PHONY: setup searxng start-searxng stop-searxng test cli app tail
###