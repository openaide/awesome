#
# make awesome
#
.DEFAULT_GOAL := help

.PHONY: help
help: Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

###
start: litellm-up postgres-up ## Start litellm services
stop: litellm-down postgres-down ## Stop litellm services

start-core: litellm-up traefik-up redis-up postgres-up ## Start litellm, traefik services
stop-core: litellm-down traefik-down redis-down postgres-down ## Stop litellm, traefik services

start-all:  start-core start-llm start-tools ## Start all services
stop-all:  stop-core stop-llm stop-tools ## Stop all services

##
.PHONY: start stop start-core stop-core start-all stop-all

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
postgres-up: ## Start postgres
	@cd perpetis && docker compose --profile postgres up -d	

postgres-down: ## Stop postgres
	@cd perpetis && docker compose --profile postgres down

redis-up: ## Start redis
	@cd perpetis && docker compose --profile redis up -d

redis-down: ## Stop redis
	@cd perpetis && docker compose --profile redis down

##
.PHONY: traefik-up traefik-down
.PHONY: litellm-up litellm-down
.PHONY: postgres-up postgres-down
.PHONY: redis-up redis-down

### LLM
ollama-up:
	@cd docker/ollama && docker compose up -d

ollama-down:
	@cd docker/ollama && docker compose down

##
start-llm: ollama-up ## Start LLM
stop-llm: ollama-down ## Stop LLM

##
.PHONY: ollama-up
.PHONY: ollama-down

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

##
start-tools: aider-up anythingllm-up docsgpt-up nextchat-up openhands-up openwebui-up ## Start tools
stop-tools: aider-down anythingllm-down docsgpt-down nextchat-down openhands-down openwebui-down ## Stop tools

##
.PHONY: aider-up anythingllm-up docsgpt-up nextchat-up openhands-up openwebui-up
.PHONY: adier-down anythingllm-down docsgpt-down nextchat-down openhands-down openwebui-down
.PHONY: start-tools stop-tools
###
