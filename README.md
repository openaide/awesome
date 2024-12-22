# Awesome AI Tools

This is a collection of open-source AI tools that anyone can run and customize on a local machine, aimed at enhancing productivity in your day-to-day workflow as a coder.

All tools are available under permissive licenses, including MIT, BSD, Apache, or MPL.

If you want to be on the cutting edge, read on.

## Requirements

* Get [Docker](https://docs.docker.com/get-started/get-docker/)

* Setting up [Visual Studio Code](https://code.visualstudio.com/docs/setup/setup-overview)

* Optional [Docker for Visual Studio Code extension](https://marketplace.visualstudio.com/items?itemName=ms-azuretools.vscode-docker)

## Getting started

```bash
git clone https://github.com/openaide/awesome.git
git submodule update --init --recursive
```

### Start gateway [LiteLLM](https://docs.litellm.ai/docs/simple_proxy) and [Traefik](https://doc.traefik.io/traefik/)

Running the gateway is optional but recommended.

```bash
cd awesome/

#
export OPENAI_API_KEY=

# you may change the default ports:
#
# LiteLLM
# export LLM_PORT=4000
#
# Traefik
# export WEB_PORT=80
# export ADMIN_PORT=8080

make start
make stop
```

Check Traefik dashboard [http://localhost:8080](http://localhost:8080)

### LLM API configuration for AI tools and VSCode extensions

> [!TIP]
>
> Add the following to /etc/hosts on your host so you could use `host.docker.internal` both in container and on host.
>
>`127.0.0.1 host.docker.internal`
>

```text
Base URL: http://<hostname>:4000
API Key: sk-1234
Model: gpt-4o #gpt-4o-mini and others
```

where `<hosthame>` is `localhost` if the tool runs on host (Continue, Cline...) or `host.docker.internal` inside docker container (OpenHands, Aider...)

The above configuration is the default assuming `OPENAI_API_KEY` env is set.
You can setup other providers with LiteLLM Admin Panel [http://localhost:4000/ui/](http://localhost:4000/ui/)

```text
Username: admin
Password: sk-1234
```

### Services

All tools can be built, started, or stopped with `docker compose`. For convenience, make targets are also provided which also take care of required dependencies.

```bash
cd docker/<app>

# docker compose build
make build

# docker compose up -d
make up

# docker compose down
make down
```

Visit the tool's web app: `http://<app>.localhost`, where `<app>` is name of the tool.
See [RFC 6761](https://www.rfc-editor.org/rfc/rfc6761) Special-Use Domain Names - 6.3.  Domain Name Reservation Considerations for `localhost`.

### VSCode extensions

All VSCode extensions can be built with make and extension is saved in `local/extension/`

```bash
make vsce
```

To install the extension: Activity Bar/Extensions/Install from VSIX...

## Tools and VSCode extensions

The following tools are grouped by their main features. The AI landscape is changing daily, the information could be inaccurate by the time you get here.

### VSCode extensions - code edit

* [Continue](https://github.com/continuedev/continue)
* [Cline](https://github.com/cline/cline.git)
* [Aider Composer](https://github.com/lee88688/aider-composer.git)

### Code generation

* [OpenHands](https://docs.all-hands.dev/)
* [Aider](docker/aider/READEME.md)
* [screenshot-to-code](docker/screenshot-to-code/REAMDE.md)
* [Bolt.diy](docker/bolt.diy/README.md)

### SQL

* [Vanna](docker/vanna/README.md)

### Document

* [AnythingLLM](https://docs.anythingllm.com/)
* [DocsGPT](https://docs.docsgpt.cloud/)
* [GPT Researcher](docker/gpt-researcher/README.md)

### Search

* [Danswer](docker/danswer/README.md)

### General

* [NextChat](https://github.com/ChatGPTNextWeb/ChatGPT-Next-Web)
* [Open WebUI](https://docs.openwebui.com/)

### LLM

* [Ollama](https://github.com/ollama/ollama?tab=readme-ov-file)
* [LocalAI](https://github.com/mudler/LocalAI?tab=readme-ov-file)

### Misc

* [MarkItDown](https://github.com/microsoft/markitdown)

### Dependencies (optional)

* [LiteLLM](https://github.com/BerriAI/litellm) proxy server (LLM gateway)
* [Traefik](https://github.com/traefik/traefik/) HTTP reverse proxy

## How to prompt

Although written for a specific tool, they could apply to all in general.

* ChatGPT [prompt engineering](https://platform.openai.com/docs/guides/prompt-engineering)
* OpenHands [prompting best practices](https://docs.all-hands.dev/modules/usage/prompting-best-practices)
* Aider [usage tips](https://aider.chat/docs/usage/tips.html)

## Benchmark

[SWE-bench Leaderboard](https://www.swebench.com/)