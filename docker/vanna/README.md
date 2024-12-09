# Vanna SQL Copilot

##

```bash
# https://github.com/vanna-ai/vanna.git
make clone
```

## Install dependencies

```bash
#
./build-app.sh
# python3 -m venv local/venv
# source local/venv/bin/activate

# pip install --upgrade pip
# pip install chromadb openai

# pip install local/vanna/[all]

#
export HOST="0.0.0.0"
export PORT="5000"

export LLM_API_KEY="sk-1234"
export LLM_BASE_URL="http://localhost:4000"
export LLM_MODEL="gpt-4o"

export POSTGRES_HOST="localhost"
export POSTGRES_PORT="5432"
export POSTGRES_DBNAME="postgres"
export POSTGRES_USER=
export POSTGRES_PASSWORD=

export STORE_BASE=./local/store
export TRAIN_BASE=./local/train

export LOG_LEVEL=debug
export FLASK_DEBUG=1

python app.py
```
