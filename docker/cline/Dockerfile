# syntax=docker/dockerfile:1

ARG NODE_VERSION=20
ARG VSCE_VERSION=3.2.1

###
FROM node:${NODE_VERSION} AS builder

WORKDIR /app

COPY ./ ./
RUN npm run install:all
RUN npm install -g @vscode/vsce@${VSCE_VERSION}

ARG VSCE_TARGET

RUN npm run vscode:prepublish
RUN vsce package --out /app/cline.vsix --target $VSCE_TARGET

###
FROM scratch AS vsce

COPY --link --from=builder /app/cline.vsix /cline.vsix
###