#!/bin/bash

set -e

# needs to be run from the root of the repository
cd local/openhands || exit 1

# create venv and install the app dependencies
python3.12 -m venv .venv

source .venv/bin/activate

pip install --upgrade --no-cache-dir pip
pip install opencv-python==4.10.0.84

make build
##