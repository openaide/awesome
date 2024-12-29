#!/bin/bash

set -e

# create venv and install the app dependencies
python3.12 -m venv .venv

# source .venv/bin/activate
export PATH=".venv/bin:$PATH"

python -m pip install --upgrade --no-cache-dir pip
python -m pip install --no-cache-dir -e "local/shell_gpt[dev,test]"
#
python -m pip install --no-cache "local/shell_gpt[litellm]"
##