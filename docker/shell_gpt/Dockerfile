# syntax=docker/dockerfile:1

###
FROM python:3-slim

ENV DEFAULT_MODEL=
ENV API_BASE_URL=http://host.docker.internal:4000
ENV USE_LITELLM=true
ENV OPENAI_API_KEY=
ENV SHELL_INTERACTION=false
ENV PRETTIFY_MARKDOWN=false
ENV OS_NAME=
ENV SHELL_NAME=bash

WORKDIR /app
COPY . /app

RUN apt-get update && apt-get install -y gcc
RUN pip install --no-cache /app[litellm] && mkdir -p /tmp/shell_gpt

VOLUME /tmp/shell_gpt

ENTRYPOINT ["sgpt"]
