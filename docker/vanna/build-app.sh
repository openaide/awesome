#!/bin/bash

set -e

# create venv and install the app dependencies

python3.12 -m venv .venv

# shellcheck disable=SC1091
source .venv/bin/activate

pip install --upgrade pip
pip install chromadb openai

pip install "local/vanna/[all]"

deactivate
##