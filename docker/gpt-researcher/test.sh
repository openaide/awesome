#!/bin/bash

# load env
source env.sh

export PATH=".venv/bin:$PATH"

python test_llm.py
python test_retriever.py
