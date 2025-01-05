#!/bin/bash

set -e

# create venv and install the app dependencies
python3.12 -m venv .venv

# source .venv/bin/activate
export PATH=".venv/bin:$PATH"

python -m pip install --upgrade --no-cache-dir pip
python -m pip install --no-cache-dir -e "local/pr-agent"

# export CONFIG__CLI_MODE = True
export OPENAI__KEY=${OPENAI_API_KEY}
export GITHUB__USER_TOKEN=${GITHUB_TOKEN}
# export GITHUB__BASE_URL=https://github.mycompany.com/api/v3

# https://qodo-merge-docs.qodo.ai/tools/
# python3 -m pr_agent.cli --pr_url <pr_url> describe
# python3 -m pr_agent.cli --pr_url <pr_url> review
# python3 -m pr_agent.cli --pr_url <pr_url> improve
# python3 -m pr_agent.cli --pr_url <pr_url> ask <your question>
# python3 -m pr_agent.cli --pr_url <pr_url> update_changelog
## # python3 -m pr_agent.cli --pr_url <pr_url> add_docs
## python3 -m pr_agent.cli --pr_url <pr_url> generate_labels

# python3 -m pr_agent.cli --issue_url <issue_url> similar_issue

# Supported commands:
# - review / review_pr - Add a review that includes a summary of the PR and specific suggestions for improvement.
# - ask / ask_question [question] - Ask a question about the PR.
# - describe / describe_pr - Modify the PR title and description based on the PR's contents.
# - improve / improve_code - Suggest improvements to the code in the PR as pull request comments ready to commit.
# Extended mode ('improve --extended') employs several calls, and provides a more thorough feedback
# - reflect - Ask the PR author questions about the PR.
# - update_changelog - Update the changelog based on the PR's contents.
# - add_docs
# - generate_labels

#
# [config]
# # models
# model="gpt-4o"
# model_weak="gpt-4o-mini"
# # CLI
# git_provider="local"
# publish_output=false
# publish_output_progress=true
# publish_output_no_suggestions=true
# # 0,1,2
# verbosity_level=2
# # Configurations
# use_wiki_settings_file=false
# use_repo_settings_file=false
# use_global_settings_file=false
# disable_auto_feedback=true
##