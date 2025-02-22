# syntax=docker/dockerfile:1

FROM rust:latest

RUN cargo install cross --locked

RUN apt-get update && apt-get install -y \
    gcc-10-multilib g++-10-multilib \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY . .

# The final binary will be located in target/x86_64-apple-darwin/release/goose
RUN cross build --target x86_64-apple-darwin --release

CMD [ "bash" ]