#
# make traefik
#

###
.DEFAULT_GOAL := help

.PHONY: help
help: Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

##
APP_VERSION ?= v3.2.1
APP_REPO_URL ?= https://github.com/traefik/traefik.git
TARGET_DIR = local/traefik

export ARCH := $(shell uname -m)

##
clone:
	@if [ ! -d $(TARGET_DIR) ]; then \
		echo "Cloning repository..."; \
		git clone --recurse-submodules --branch $(APP_VERSION) $(APP_REPO_URL) $(TARGET_DIR); \
	else \
		echo "Directory exists."; \
	fi; \
	cd $(TARGET_DIR) && git pull origin $(APP_VERSION)

##
build: clone ## Build
	@docker compose build

up: net ## Start
	@docker compose up -d

down:  ## Stop
	@docker compose down

net: ## Create openland network
	@docker network inspect openland >/dev/null 2>&1 \
		|| docker network create openland

.PHONY: clone build up down net
###
