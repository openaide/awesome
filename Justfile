#!/usr/bin/env -S just --justfile

default:
    @just --list

# build the gateway using Docker
gateway-build:
    cd docker/gateway && docker buildx bake --progress=plain --file ./compose.override.yml gw

# bring up the gateway using Docker Compose
gateway-up:
    cd docker/gateway && docker compose up -d

# bring down the gateway using Docker Compose
gateway-down:
    cd docker/gateway && docker compose down

# create a Docker network if it doesn't exist
net:
    docker network inspect openland >/dev/null 2>&1 || docker network create openland

# see AI-based tool for generating commit messages
# https://github.com/qiangli/ai

# generate a git commit message and copy it to clipboard
git-message:
    git diff origin/main | ai @git/conventional }

# generate a git commit message, copy it to clipboard, and commit
git-commit: git-message
    git commit -m "$(pbpaste)"

# generate a git commit message, copy it to clipboard, and amend the commit
git-amend: git-message
    git commit --amend -m "$(pbpaste)"

# Build MCP servers - stargate
mcp-build:
    cd docker/mcp-servers/stargate && docker buildx bake --progress=plain --file ./compose.yml --file ./compose.override.yml stargate

# start the MCP servers - stargate
mcp-up:
    cd docker/mcp-servers/stargate && docker compose up -d

# stop the MCP servers - stargate
mcp-down:
    cd docker/mcp-servers/stargate && docker compose down