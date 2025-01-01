#!/bin/bash

# load env
source env.sh

VIRTUAL_ENV=".venv"

export PATH="$VIRTUAL_ENV/bin:$PATH"

python app.py
