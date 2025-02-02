#
# make awesome
#
.DEFAULT_GOAL := help

.PHONY: help
help: Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

##
export COMPOSE_PROJECT_NAME = awesome
# export COMPOSE_REMOVE_ORPHANS = true
export COMPOSE_IGNORE_ORPHANS = true

###
start: net gateway-up ## Start gateway (one litellm/traefik proxy)
stop: gateway-down ## Stop gateway (one litellm/traefik proxy)

##
.PHONY: start stop

###

##
gateway-build:
	@cd docker/gateway && docker buildx bake --progress=plain --file ./compose.override.yml gw

gateway-up:
	@cd docker/gateway && docker compose up -d

gateway-down:
	@cd docker/gateway && docker compose down

.PHONY: gateway-up gateway-down gateway-build

##
net: ## Create network
	@docker network inspect openland >/dev/null 2>&1 \
		|| docker network create openland
.PHONY: net

# https://github.com/qiangli/ai
git-message: ## Generate commit message and copy the message to clipboard
	@git diff origin/main|ai @git/conventional =+

git-commit: git-message ## Generate and commit with the message
	@git commit -m "$$(pbpaste)"

git-amend: git-message ## Generate and amend with the message
	@git commit --amend -m "$$(pbpaste)"

.PHONY: commit-message
###
