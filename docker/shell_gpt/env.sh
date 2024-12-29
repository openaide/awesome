#!/bin/bash

export DEFAULT_MODEL=gpt-4o-mini
export API_BASE_URL=http://localhost:4000
export USE_LITELLM=true
export OPENAI_API_KEY=sk-1234
export SHELL_INTERACTION=true
export PRETTIFY_MARKDOWN=false
# shellcheck disable=SC2155
export OS_NAME=$(uname -s)
export SHELL_NAME=bash
