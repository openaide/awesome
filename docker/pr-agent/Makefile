#
# make pr agent
#

###
.DEFAULT_GOAL := help

.PHONY: help
help: Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

##
APP_VERSION ?= main
APP_REPO_URL ?= https://github.com/Codium-ai/pr-agent.git
TARGET_DIR = local/pr-agent

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
build: clone checkout ## Build
	@docker compose build

run:  ## Run
	@docker compose run --rm pr-agent

.PHONY: build run

setup: ## setup venv (python)
	@setup.sh
.PHONY: setup
###
