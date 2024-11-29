#!/bin/bash
set -xeuo pipefail

function source_env() {
  if [ -f .env ]; then
    export $(cat .env | grep -v '#' | xargs)
  fi
}

##
source_env

export WEB_PORT="${WEB_PORT:-80}"
export ADMIN_PORT="${ADMIN_PORT:-8080}"

##
DOMAIN="localhost"
URLS=(
"http://${DOMAIN}:${ADMIN_PORT}/"
#
"http://gateway.${DOMAIN}:${WEB_PORT}/"
"http://litellm.${DOMAIN}:${WEB_PORT}/"
# tool service must be running
"http://aider.${DOMAIN}:${WEB_PORT}/"
"http://anythingllm.${DOMAIN}:${WEB_PORT}/"
"http://docsgp.${DOMAIN}:${WEB_PORT}/"
"http://nextchat.${DOMAIN}:${WEB_PORT}/"
"http://openhands.${DOMAIN}:${WEB_PORT}/"
"http://openwebui.${DOMAIN}:${WEB_PORT}/"
)

DATA_DIR="openaide"

##
function chrome() {
  case "$OSTYPE" in
    linux*)
      google-chrome "$@" --user-data-dir=/tmp/"$DATA_DIR" "${URLS[@]}"
      ;;
    darwin*)
      open -n -a "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" --args "$@" --user-data-dir=/tmp/"$DATA_DIR" "${URLS[@]}"
      ;;
    msys*)
      "C:\Program Files (x86)\Google\Chrome\Application\chrome.exe" "$@" --disable-gpu --user-data-dir=~/temp/"$DATA_DIR" "${URLS[@]}"
      ;;
    *)
      echo "$OSTYPE not supported"
      ;;
  esac
}

chrome "$@"
##