# syntax=docker/dockerfile:1

# ###
# # git clone --branch go-cli https://github.com/qiangli/OpenHands.git local/OpenHands
# FROM scratch AS source
# COPY ./local/OpenHands /app

###
FROM ubuntu:22.04 AS source

ARG VERSION

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    ca-certificates \
    curl \
    git && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

ARG CLI_VERSION=go-cli
ARG CLI_REPO_URL=https://github.com/qiangli/OpenHands.git

RUN git clone --recurse-submodules --branch ${CLI_VERSION}  ${CLI_REPO_URL} /app && \
    cd /app/cli && \
    git pull origin ${CLI_VERSION}

###
FROM golang:1.23.1-alpine AS build
RUN apk add --no-cache bash

ARG VERSION
ENV VERSION=${VERSION}

ARG SANDBOX
ARG IMAGE

WORKDIR /app

COPY --from=source /app/cli/go.mod /app/cli/go.sum ./
RUN go mod download

COPY --from=source /app/cli/ /app/
RUN ./build-cli.sh

###
FROM scratch
COPY --from=build /app/dist /

CMD [ "bash" ]
