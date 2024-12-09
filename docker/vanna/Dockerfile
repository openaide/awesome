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
FROM golang:1.22.3 AS proxy-builder

ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -a -installsuffix cgo -o bin/proxy ./

###
FROM python:3.12.3-slim AS app

RUN apt-get update && apt-get install -y bash && rm -rf /var/lib/apt/lists/*

WORKDIR /

COPY --from=proxy-builder /app/bin/proxy /
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
ENV POSTGRES_USER="postgres"
ENV POSTGRES_PASSWORD=""

ENV STORE_BASE="/workspace/store"
ENV TRAIN_BASE="/workspace/train"

ENV STORE_PATH="${STORE_BASE}/${POSTGRES_USER}/${POSTGRES_DBNAME}"
ENV TRAIN_PATH="${TRAIN_BASE}"

ENV TOKENIZERS_PARALLELISM=false
ENV LOG_LEVEL=debug
ENV FLASK_DEBUG=1

VOLUME [ "/workspace" ]

CMD [ "/proxy", "--address", "0.0.0.0:58080", "--venv", "/venv", "--script", "/app.py" ]
###
