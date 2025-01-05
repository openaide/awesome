#!/bin/bash

set -e

# create venv and install the app dependencies

python3.12 -m venv .venv

#
# shellcheck disable=SC1091
source .venv/bin/activate

pip install --upgrade --no-cache-dir pip
pip install --no-cache-dir chromadb openai
pip install -e "local/vanna/[all]"

# deactivate
##