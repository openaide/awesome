#!/bin/bash

# query=$1

# # research_report
# # detailed_report
# # outline_report
# report_type=${2:-"research_report"}

# query="What are the main causes of climate change?"
# report_type="research_report"

#
query="The impact of artificial intelligence on job markets"
report_type="research_report"
# report_type="detailed_report"

# #
# query="Renewable energy sources and their potential"
# report_type="outline_report"


##
IMAGE="openaide/gptr-cli"

docker build -t "$IMAGE" -f Dockerfile.cli .

docker run --rm --env-file "env_vars.txt" --env OPENAI_API_BASE="http://host.docker.internal:4000/" --volume ./outputs/:/app/outputs/ "$IMAGE" "$query" --report_type "$report_type"
