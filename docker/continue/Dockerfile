# syntax=docker/dockerfile:1

ARG NODE_VERSION=20.19.0

###
FROM node:${NODE_VERSION} AS builder

WORKDIR /app

COPY ./package*.json ./
RUN npm install

#
COPY ./ ./

ARG VSCE_TARGET

# TODO no all dependencies are required for vscode extension
RUN bash scripts/install-dependencies.sh

# .github/workflows/main.yml
RUN <<EOF
#!/bin/bash
set -xeuo pipefail

touch /.dockerenv

cd /app/core && npm install && npm i vectordb
cd /app/gui && npm install && npm run build

cd /app/extensions/vscode || exit 1
rm -rf build/

npm install
npm install -g @vscode/vsce

npm run prepackage -- --target $VSCE_TARGET
npm install -f esbuild
npx vsce package --out /app/continue.vsix --no-dependencies --target $VSCE_TARGET
EOF

###
FROM scratch AS vsce

COPY --link --from=builder /app/continue.vsix /
###