#!/bin/bash

set -e

# create venv and install the app dependencies
python3.12 -m venv .venv

source .venv/bin/activate
export PATH=".venv/bin:$PATH"

pip install --upgrade --no-cache-dir pip

# local/chatbot/requirements.txt
pip install -r ./requirements.txt

##