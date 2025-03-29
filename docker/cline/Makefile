#
# make cline
#

###
.DEFAULT_GOAL := help

.PHONY: help
help: Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

##
APP_VERSION ?= main
# APP_VERSION ?= v3.8.4
APP_REPO_URL ?= https://github.com/cline/cline.git
TARGET_DIR = local/cline

PATCH = $(PWD)/src.patch
TIMESTAMP := $(shell date +'%Y-%m-%d %H:%M:%S')

OS := $(shell uname)
ARCH := $(shell uname -m)

VSCE_TARGET := $(if $(findstring Darwin,$(OS)),darwin,linux)-$(if $(findstring arm64,$(ARCH)),arm64,x64)

info: ## Show info
	@echo "VSCE_TARGET: $(VSCE_TARGET)"
	@echo "APP_REPO_URL: $(APP_REPO_URL)"
	@echo "APP_VERSION: $(APP_VERSION)"
	@echo "TARGET_DIR: $(TARGET_DIR)"
	@echo "VSCode Extension: ./local/extension"

.PHONY: info

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

patch:
	@cd $(TARGET_DIR) && git reset --hard && git apply $(PATCH)
	sed "s/\[\[BUILD_TIMESTAMP\]\]/$(TIMESTAMP)/g" $(TARGET_DIR)/webview-ui/src/components/chat/ChatView.tsx > local/ChatView.tsx && \
	mv local/ChatView.tsx $(TARGET_DIR)/webview-ui/src/components/chat/ChatView.tsx

.PHONY: clone checkout patch timestamp

##
build: clone checkout patch vsce
vsce: clean ## Build cline vscode extension
	@docker buildx build --build-arg VSCE_TARGET=$(VSCE_TARGET) --progress plain --target vsce --output type=local,dest=./local/extension --file ./Dockerfile $(TARGET_DIR)

clean:
	@rm -rf ./local/extension

.PHONY: build vsce clean
###
