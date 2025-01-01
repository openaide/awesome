#!/bin/bash

# load env
source env.sh

export PATH=".venv/bin:$PATH"

query=${1:-"What are the main causes of climate change?"}
report_type=${2:-"research_report"}

python local/gpt-researcher/cli.py "$query"  --report_type "$report_type" 
