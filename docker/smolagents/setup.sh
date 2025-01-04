#!/bin/bash

set -e

# create venv and install the app dependencies
python3.12 -m venv .venv

# source .venv/bin/activate
VIRTUAL_ENV=".venv"
export PATH="$VIRTUAL_ENV/bin:$PATH"

python -m pip install --upgrade --no-cache-dir pip
python -m pip install --no-cache-dir -r "./requirements.txt"

python -m pip install --no-cache-dir -e "local/smolagents/"

##