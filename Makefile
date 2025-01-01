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
build: gateway-build ## Build gateway (one litellm/traefik proxy)

start: net gateway-up ## Start gateway (one litellm/traefik proxy)
stop: gateway-down ## Stop gateway (one litellm/traefik proxy)

start-all: start-gw start-llm start-tools ## Start all services
stop-all: stop-gw stop-llm stop-tools ## Stop all services

start-proxy: net litellm-up postgres-up traefik-up ## Start litellm and traefik proxy services
stop-proxy: litellm-down postgres-down traefik-down ## Stop litellm and traefik proxy services

##
.PHONY: build start stop start-all stop-all start-proxy stop-proxy

###
traefik-up: ## Start traefik
	@cd docker/traefik && docker compose up -d	

traefik-down: ## Stop traefik
	@cd docker/traefik && docker compose down

#
litellm-up: ## Start litellm
	@cd docker/litellm && docker compose up -d

litellm-down: ## Stop litellm
	@cd docker/litellm && docker compose down

##
gateway-build:
	@cd docker/gateway && docker buildx bake --progress=plain --file ./compose.override.yml gw

gateway-up:
	@cd docker/gateway && docker compose up -d

gateway-down:
	@cd docker/gateway && docker compose down

.PHONY: gateway-up gateway-down gateway-build

##
postgres-up:
	@cd perpetis && docker compose --profile postgres up -d	

postgres-down:
	@cd perpetis && docker compose --profile postgres down

redis-up:
	@cd perpetis && docker compose --profile redis up -d

redis-down:
	@cd perpetis && docker compose --profile redis down

##
.PHONY: traefik-up traefik-down
.PHONY: litellm-up litellm-down
.PHONY: postgres-up postgres-down
.PHONY: redis-up redis-down

### LLM
ollama-up: ## Start ollama
	@cd docker/ollama && docker compose up -d

ollama-down: ## Stop ollama
	@cd docker/ollama && docker compose down

localai-up:
	@cd docker/localai && docker compose up -d

localai-down:
	@cd docker/localai && docker compose down

##
start-llm: ollama-up localai-up ## Start all LLM
stop-llm: ollama-down localai-down ## Stop all LLM

##
.PHONY: ollama-up localai-up
.PHONY: ollama-down localai-down

.PHONY: start-llm stop-llm

### tool apps
aider-up:
	@cd docker/aider && docker compose up -d

aider-down:
	@cd docker/aider && docker compose down

#
anythingllm-up:
	@cd docker/anythingllm && docker compose up -d 

anythingllm-down:
	@cd docker/anythingllm && docker compose down

#
docsgpt-up:
	@cd docker/docsgpt && docker compose up -d 

docsgpt-down:
	@cd docker/docsgpt && docker compose down

#
nextchat-up:
	@cd docker/nextchat && docker compose up -d

nextchat-down:
	@cd docker/nextchat && docker compose down

#
openhands-up:
	@cd docker/openhands && docker compose up -d 

openhands-down:
	@cd docker/openhands && docker compose down

#
openwebui-up: ## Start open webui
	@cd docker/openwebui && docker compose up -d

openwebui-down: ## Stop open webui
	@cd docker/openwebui && docker compose down

vanna-up: ## Start vanna
	@cd docker/vanna && docker compose up -d

vanna-down: ## Stop vanna
	@cd docker/vanna && docker compose down

##
start-tools: aider-up anythingllm-up docsgpt-up nextchat-up openhands-up openwebui-up ## Start tools
stop-tools: aider-down anythingllm-down docsgpt-down nextchat-down openhands-down openwebui-down ## Stop tools

##
.PHONY: aider-up anythingllm-up docsgpt-up nextchat-up openhands-up openwebui-up
.PHONY: adier-down anythingllm-down docsgpt-down nextchat-down openhands-down openwebui-down
.PHONY: start-tools stop-tools

.PHONY: vanna-up vanna-down

##
net: ## Create network
	@docker network inspect openland >/dev/null 2>&1 \
		|| docker network create openland
.PHONY: net

# https://github.com/qiangli/ai
commit-message: ## Generate commit message and copy the message to clipboard
	@git diff origin main|ai @ask write commit message for git =+

.PHONY: commit-message
###
