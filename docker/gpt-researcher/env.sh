#!/bin/bash

# shellcheck disable=all

# https://docs.gptr.dev/docs/gpt-researcher/gptr/config
# https://github.com/assafelovic/gpt-researcher/blob/master/docs/docs/gpt-researcher/gptr/config.md

export OPENAI_API_BASE=http://localhost:4000/v1
export OPENAI_API_KEY=sk-1234
# export RETRIEVER=google
# export GOOGLE_API_KEY=""

# export RETRIEVER=custom
# # export RETRIEVER_ARG_API_KEY=
# export RETRIEVER_ENDPOINT=http://localhost:48080/

export RETRIEVER=searx
export SEARX_URL=http://localhost:48080/
# export RETRIEVER_ENDPOINT=http://localhost:48080/

#
export EMBEDDING=openai:text-embedding-3-small
export FAST_LLM=openai:gpt-4o-mini
export SMART_LLM=openai:gpt-4o
export STRATEGIC_LLM=openai:o1-preview

# export EMBEDDING=ollama:nomic-embed-text
# export OLLAMA_BASE_URL=http://ollama:11434
# export FAST_LLM=ollama:llama3
# export SMART_LLM=ollama:llama3
# export STRATEGIC_LLM=ollama:llama3
#
export CURATE_SOURCES=True
export REPORT_FORMAT=IEEE
# export REPORT_FORMAT=Infographic
export DOC_PATH=$PWD/inputs/docs
#
export SCRAPER=bs
# export SCRAPER=browser