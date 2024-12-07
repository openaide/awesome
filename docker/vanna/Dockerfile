# syntax=docker/dockerfile:1

###
FROM python:3.12.3-slim AS builder

RUN apt-get update -y && apt-get install -y \
	bash \
	curl \
	git

WORKDIR /app

COPY ./local/vanna ./

RUN python3 -m venv /venv

RUN /venv/bin/pip install --upgrade pip \
    && /venv/bin/pip install chromadb openai \
    && /venv/bin/pip install "./[all]"

###
FROM python:3.12.3-slim AS app

COPY --from=builder /venv /venv
COPY app.py /app.py

ENV HOST="0.0.0.0"
ENV PORT="5000"

ENV LLM_API_KEY=
ENV LLM_BASE_URL="http://host.docker.internal:4000"
ENV LLM_MODEL="gpt-4o"

ENV POSTGRES_HOST="host.docker.internal"
ENV POSTGRES_PORT="5432"
ENV POSTGRES_DBNAME="postgres"
ENV POSTGRES_USER=""
ENV POSTGRES_PASSWORD=""

ENV STORE_PATH="/workspace"

VOLUME [ "/workspace" ]

CMD [ "/venv/bin/python", "/app.py" ]