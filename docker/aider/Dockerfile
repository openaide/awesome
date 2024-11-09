# syntax=docker/dockerfile:1

###
# git clone --branch main https://github.com/Aider-AI/aider.git local/aider
FROM scratch AS source
COPY ./local/aider /app

###
FROM ubuntu:22.04 AS sourcex

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    ca-certificates \
    curl \
    git && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

ARG AIDER_VERSION=v0.62.1

RUN git clone --branch ${AIDER_VERSION}  https://github.com/Aider-AI/aider.git /app

###
FROM python:3.10-slim AS base

RUN apt-get update && \
    apt-get install --no-install-recommends -y build-essential git libportaudio2 pandoc && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

RUN python -m venv /venv
ENV PATH="/venv/bin:$PATH"

# https://playwright.dev/python/docs/browsers
ENV PLAYWRIGHT_BROWSERS_PATH=/pw-browsers
ENV PLAYWRIGHT_SKIP_BROWSER_GC=1

# Permission kludges to support `docker run --user xxx`
RUN chmod a+rwx /venv /venv/bin /venv/include /venv/lib /venv/lib/python3.10/site-packages

RUN mkdir       /.aider /.cache /pw-browsers
RUN chmod a+rwx /.aider /.cache /pw-browsers

# So git doesn't complain about unusual permissions
RUN git config --system --add safe.directory /app

#########################
FROM base AS aider-full

ENV AIDER_DOCKER_IMAGE=paulgauthier/aider-full

COPY --from=source /app /tmp/aider

RUN /venv/bin/python -m pip install --upgrade --no-cache-dir pip \
    && /venv/bin/python -m pip install --no-cache-dir /tmp/aider[help,browser,playwright] \
       --extra-index-url https://download.pytorch.org/whl/cpu \
    && rm -rf /tmp/aider

RUN /venv/bin/python -m playwright install --with-deps chromium
RUN find /venv/lib/python3.10/site-packages \( -type d -exec chmod a+rwx {} + \) -o \( -type f -exec chmod a+rw {} + \)

ENTRYPOINT ["/venv/bin/aider"]

#########################
FROM base AS aider

ENV AIDER_DOCKER_IMAGE=paulgauthier/aider

COPY --from=source /app /tmp/aider

RUN /venv/bin/python -m pip install --upgrade --no-cache-dir pip \
    && /venv/bin/python -m pip install --no-cache-dir /tmp/aider[playwright] \
       --extra-index-url https://download.pytorch.org/whl/cpu \
    && rm -rf /tmp/aider

RUN /venv/bin/python -m playwright install --with-deps chromium
RUN find /venv/lib/python3.10/site-packages \( -type d -exec chmod a+rwx {} + \) -o \( -type f -exec chmod a+rw {} + \)

ENTRYPOINT ["/venv/bin/aider"]

#########################
# https://aider.chat/docs/faq.html#how-can-i-run-aider-locally-from-source-code
# https://aider.chat/docs/llms/openai-compat.html
FROM base AS aider-source

COPY --from=source /app /aider

WORKDIR /aider
RUN python -m pip install -e .

ENV PIP_NO_INPUT=1
RUN python -m pip install --upgrade --upgrade-strategy only-if-needed 'aider-chat[browser]'

ENV OPENAI_API_BASE="http://host.docker.internal:4000"
ENV OPENAI_API_KEY="sk-1234"

CMD [ "python", "-m", "aider", "--gui", "--model", "openai/gpt-4o-mini" ]
###
