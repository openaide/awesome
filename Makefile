#
# make awesome
#
.DEFAULT_GOAL := help

.PHONY: help
help: Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

##
start: litellm-up postgres-up ## Start api key proxy services
stop: litellm-down postgres-down ## Stop api key proxy services

start-all: start traefik-up redis-up ## Start all services
stop-all: stop traefik-down redis-down ## Stop all services

##
traefik-up: ## Start traefik
	@cd docker/traefik && docker compose up -d	

traefik-down: ## Stop traefik
	@cd docker/traefik && docker compose down

##
postgres-up: ## Start postgres
	@cd perpetis && docker compose --profile postgres up -d	

postgres-down: ## Stop postgres
	@cd perpetis && docker compose --profile postgres down

redis-up: ## Start redis
	@cd perpetis && docker compose --profile redis up -d

redis-down: ## Stop redis
	@cd perpetis && docker compose --profile redis down

##
litellm-up: ## Start litellm
	@cd docker/litellm && docker compose up -d

litellm-down: ## Stop litellm
	@cd docker/litellm && docker compose down

##
.PHONY: start stop start-all stop-all
.PHONY: traefik-up traefik-down
.PHONY: litellm-up litellm-down
.PHONY: postgres-up postgres-down redis-up redis-down
.PHONY: tabby
##
