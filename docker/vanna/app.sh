#!/bin/bash

set -e

##
export HOST="0.0.0.0"
export PORT="${PORT:-5000}"

export LLM_API_KEY="sk-1234"
export LLM_BASE_URL="http://localhost:4000"
export LLM_MODEL="gpt-4o"

export POSTGRES_HOST="localhost"
export POSTGRES_PORT="${POSTGRES_PORT:-5432}"
export POSTGRES_DBNAME=${POSTGRES_DBNAME:-"postgres"}
export POSTGRES_USER=${POSTGRES_USER:-"postgres"}
export POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-}

export STORE_BASE=${STORE_BASE:-"./local/store"}
export TRAIN_BASE=${TRAIN_BASE:-"./local/train"}

#
export PATH=./local/venv/bin:$PATH

# 
export TOKENIZERS_PARALLELISM=false
export LOG_LEVEL=debug
export FLASK_DEBUG=1

python app.py "$@"
##
