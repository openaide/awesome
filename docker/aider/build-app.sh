#!/bin/bash

set -e

# create venv and install the app dependencies
python3.12 -m venv .venv

export PATH=".venv/bin:$PATH"

python -m pip install --upgrade --no-cache-dir pip
python -m pip install --no-cache-dir -e "local/aider[help,browser,playwright]"

##