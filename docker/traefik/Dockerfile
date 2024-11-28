# syntax=docker/dockerfile:1

###
FROM ubuntu:22.04 AS dind

# https://docs.docker.com/engine/install/ubuntu/
RUN apt-get update && apt-get install -y \
	ca-certificates \
	curl \
	&& install -m 0755 -d /etc/apt/keyrings \
	&& curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc \
	&& chmod a+r /etc/apt/keyrings/docker.asc \
	&& echo \
		"deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
		$(. /etc/os-release && echo "$VERSION_CODENAME") stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

RUN apt-get update && apt-get install -y \
	docker-ce \
	docker-ce-cli \
	containerd.io \
	docker-buildx-plugin \
	docker-compose-plugin \
	&& rm -rf /var/lib/apt/lists/* \
	&& apt-get clean \
	&& apt-get autoremove -y

###
FROM dind AS builder

ENV DEBIAN_FRONTEND=noninteractive

#
RUN apt-get update && apt-get install -y \
	bash \
    build-essential \
    curl \
	git \
	git-lfs \
	software-properties-common \
	make \
    netcat \
    sudo \
	tar \
	wget \
	&& rm -rf /var/lib/apt/lists/* \
	&& apt-get clean \
	&& apt-get autoremove -y

ARG ARCH
ARG GO_VERSION=1.23.3

# Install Go
RUN curl -fsSL -o go${GO_VERSION}.linux-${ARCH}.tar.gz https://go.dev/dl/go${GO_VERSION}.linux-${ARCH}.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-${ARCH}.tar.gz && \
    rm go${GO_VERSION}.linux-${ARCH}.tar.gz

#
ENV PATH="/usr/local/go/bin:${PATH}"

###
FROM builder AS traefik

WORKDIR /app

COPY ./local/traefik /app

RUN make

###
# local/traefik/Dockerfile
FROM alpine:3.20

RUN apk add --no-cache --no-progress ca-certificates tzdata

ARG TARGETPLATFORM
COPY --from=traefik /app/dist/$TARGETPLATFORM/traefik /

EXPOSE 80
VOLUME ["/tmp"]

ENTRYPOINT ["/traefik"]

###