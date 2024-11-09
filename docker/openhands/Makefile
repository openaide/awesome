#
# make openhands
#

###
.DEFAULT_GOAL := help

.PHONY: help
help: Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

##
build: ## Build image
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

.PHONY: build up down net
.PHONY: cli clean
###
