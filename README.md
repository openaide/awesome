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

### Start gateway (LiteLLM/Traefik)

Running the gateway is optional but recommended. See [LiteLLM](https://docs.litellm.ai/docs/simple_proxy) and [Traefik](https://doc.traefik.io/traefik/) for more information.

```bash
cd awesome/

make build

#
export OPENAI_API_KEY=

# you may change the defaults:
#
# LiteLLM
# export LLM_PORT=4000
# export LITELLM_SALT_KEY=sk-1234
# export LITELLM_MASTER_KEY=sk-1234
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

Configure LiteLLM models and settings here: [docker/gateway/etc/service/litellm/config.yaml](docker/gateway/etc/service/litellm/config.yaml).

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

### Command line

* [AI](https://github.com/qiangli/ai.git)
* [ShellGPT](https://github.com/TheR1D/shell_gpt)

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
* [Midscene.js](https://github.com/web-infra-dev/midscene)
* [smolagents](https://github.com/huggingface/smolagents)
* [SWE-agent](https://github.com/SWE-agent/SWE-agent)
* [CodeAct](https://github.com/All-Hands-AI/OpenHands/tree/main/openhands/agenthub/codeact_agent)
* [STORM](https://github.com/stanford-oval/storm)
* [Vector Databases](https://cookbook.openai.com/examples/vector_databases/readme)
* [Sentence Transformers](https://www.sbert.net/)
* [PR Agent](https://github.com/Codium-ai/pr-agent)
* [Bloop](https://github.com/BloopAI/bloop)
* [Awesome AI Tools](https://github.com/mahseema/awesome-ai-tools)
* [Awesome AI](https://github.com/re50urces/Awesome-AI)

### Dependencies (optional)

* [LiteLLM](https://github.com/BerriAI/litellm) proxy server (LLM gateway)
* [Traefik](https://github.com/traefik/traefik/) HTTP reverse proxy

## How to prompt

* [PromptTools](https://github.com/hegelai/prompttools)
* [Prompt Engineering Guide](https://www.promptingguide.ai/)
* [Prompting Guides](https://cookbook.openai.com/articles/related_resources#prompting-guides)
* ChatGPT [prompt engineering](https://platform.openai.com/docs/guides/prompt-engineering)
* OpenHands [prompting best practices](https://docs.all-hands.dev/modules/usage/prompting-best-practices)
* Aider [usage tips](https://aider.chat/docs/usage/tips.html)

## Benchmark

* SWE-bench [Leaderboard](https://www.swebench.com/)
* Massive Text Embedding Benchmark (MTEB) [Leaderboard](https://huggingface.co/spaces/mteb/leaderboard)