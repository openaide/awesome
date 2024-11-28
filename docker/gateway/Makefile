#
# make gateway
#

###
.DEFAULT_GOAL := help

.PHONY: help
help: Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

##
export ARCH := $(shell uname -m)

##
build: ## Build gateway
	@docker buildx bake --progress=plain --file ./compose.override.yml gw

up: net ## Start gateway
	@docker compose up -d

down: ## Stop gateway
	@docker compose down

net: ## Create openland network
	@docker network inspect openland >/dev/null 2>&1 \
		|| docker network create openland

.PHONY: build up down net
###
