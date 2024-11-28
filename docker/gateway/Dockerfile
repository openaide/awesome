# syntax=docker/dockerfile:1

###
FROM openaide/traefik AS traefik

EXPOSE 80 8080

###
FROM openaide/litellm-osi AS litellm

EXPOSE 4000

FROM litellm AS postgres

ENV POSTGRES_USER=postgres

ARG POSTGRES_VERSION=17

RUN apt-get update && \
    apt-get install -y gnupg2 wget lsb-release && \
    echo "deb http://apt.postgresql.org/pub/repos/apt/ $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list && \
    wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - && \
    apt-get update && \
    apt-get install -y postgresql-${POSTGRES_VERSION} postgresql-contrib-${POSTGRES_VERSION} && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

ENV PATH="/usr/lib/postgresql/${POSTGRES_VERSION}/bin:${PATH}"

USER postgres
RUN initdb -D /var/lib/postgresql/data

EXPOSE 5432

###
FROM postgres AS gateway

USER root

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    runit && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY --from=traefik /traefik /opt/bin/

COPY ./entrypoint.sh /entrypoint.sh
COPY ./entrypoint.d /entrypoint.d/
COPY ./opt/ /opt/
COPY ./etc/service /etc/service

ENV PATH="/opt/bin:${PATH}"

##
ENTRYPOINT ["/entrypoint.sh"]

CMD ["runsvdir", "/etc/service"]
