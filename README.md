# Awesome AI Tools

This is a collection of open-source AI tools and dependencies that anyone can run and customize on a local machine, aimed at enhancing productivity in your day-to-day workflow as a coder.

All tools are available under permissive licenses, including MIT, BSD, Apache, or MPL.

If you want to be on the cutting edge, read on.

## Requirements

* Get [Docker](https://docs.docker.com/get-started/get-docker/)

* Setting up [Visual Studio Code](https://code.visualstudio.com/docs/setup/setup-overview)

* Optional extension [Docker for Visual Studio Code](https://marketplace.visualstudio.com/items?itemName=ms-azuretools.vscode-docker)

## Getting started

```bash
git clone https://github.com/openaide/awesome.git
```

Start LLM gateway [LiteLLM](https://github.com/BerriAI/litellm).

Running the gateway is optional but recommended.

```bash
cd awesome/

make start
make stop
```

> [!TIP]
>
> Add the following to /etc/hosts on your host
>
>`127.0.0.1 host.docker.internal`
>

LLM API configuration

```text
Base URL: http://<hostname>:4000
API Key: sk-1234
Model: gpt-4o #or gpt-4o-mini
```

where `<hosthame>` is `localhost` if the tool runs on host (Continue, Cline) or `host.docker.internal` inside docker container (OpenHands, Aider...)

The above configuration is the default assuming `OPENAI_API_KEY` env is set.
You can setup other providers with LiteLLM Admin Panel [http://localhost:4000/ui/](http://localhost:4000/ui/)

```text
Username: admin
Password: sk-1234
```

### Services

All tools can be built, started, or stopped with `docker compose`. For convenience, make targets are also provided.

```bash
cd docker/<app>
make build

make up
make down
```

Start the application proxy [traefik](https://doc.traefik.io/traefik/) including the LLM gateway.

```bash
make start-all
make stop-all
```

Check traefik dashboard [http://localhost:8080](http://localhost:8080)

Visit the tool's web app: `http://<app>.localhost` where `<app>` is name of the tool.
Domain Name Reservation Considerations for [localhost](https://www.rfc-editor.org/rfc/rfc6761)

### VSCode extensions

All VSCode extensions can be built with make and extension is saved in `local/extension/`

```bash
make vsce
```

To install the extension: Activity Bar/Extensions/Install from VSIX...

## Tools and extensions

The following tools are grouped by their main features. The AI landscape is changing daily, the information could be inaccurate by the time you get here.

### VSCode extension - code edit

* [Continue](https://github.com/continuedev/continue)
* [Cline](https://github.com/cline/cline.git)
* [Aider Composer](https://github.com/lee88688/aider-composer.git)

### Code generation

* [OpenHands](https://docs.all-hands.dev/)
* [Aider](https://aider.chat/docs/usage/browser.html)
* [screenshot-to-code](docker/screenshot-to-code/REAMDE.md)

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

### Dependencies

## How to prompt

Although written for a specific tool, they could apply to all in general.

* ChatGPT [prompt engineering](https://platform.openai.com/docs/guides/prompt-engineering)
* OpenHands [prompting best practices](https://docs.all-hands.dev/modules/usage/prompting-best-practices)
* Aider [usage tips](https://aider.chat/docs/usage/tips.html)