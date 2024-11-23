#!/bin/bash

set -xeuo pipefail

OPENAIDE_BASE=$(git rev-parse --show-toplevel)

APP_DIR="${OPENAIDE_BASE}/docker/continue/local/continue"

## dependencies
cd "${APP_DIR}" || exit 1
bash scripts/install-dependencies.sh

##
cd "${APP_DIR}/extensions/vscode" || exit 1
npm install
npm run vscode:prepublish
npm run package

echo "Done building continue.dev vscode extension"
echo "Install the extension in vscode:"
echo "Activity Bar -> Extensions -> Install from VSIX..."
echo "${APP_DIR}/extensions/vscode/build/continue-*.vsix"
###