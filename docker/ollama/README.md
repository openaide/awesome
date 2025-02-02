# Ollama

## Build and run

Build could take a long time depending on your host configurations. You may skip building by directly pulling the Ollama image.

```bash
make build
make up
make down
```

Prepare models once ollama is up and running:

```bash
#
# docker exec ollama ollama pull llama3.2
# docker exec ollama ollama pull llama3.2:1b
# docker exec ollama ollama pull starcoder2:3b
# docker exec ollama ollama pull codellama:7b-instruct
# docker exec ollama ollama pull codellama:7b-code
# docker exec ollama ollama pull qwen2.5-coder:7b-base
docker exec ollama ollama pull deepseek-r1:7b
docker exec ollama ollama pull deepseek-r1:14b
#
docker exec ollama ollama list
docker exec ollama ollama show deepseek-r1:7b
```

* [Example models](https://github.com/ollama/ollama?tab=readme-ov-file)
* [Model library](https://ollama.com/library)

## Integration

### Open WebUI

[Bulid and run](../openwebui/README.md)

### Continue.Dev

[Model setup](https://docs.continue.dev/autocomplete/model-setup)

```json
"models": [
  {
    "model": "gpt-4o",
    "provider": "openai",
    "systemMessage": "You are an expert software developer. You give helpful and concise responses.",
    "apiKey": "sk-1234",
    "apiBase": "http://localhost:4000",
    "title": "GPT-4o"
  },
  {
    "model": "gpt-4o-mini",
    "provider": "openai",
    "systemMessage": "You are an expert software developer. You give helpful and concise responses.",
    "apiKey": "sk-1234",
    "apiBase": "http://localhost:4000",
    "title": "GPT-4o Mini"
  },
  {
    "model": "deepseek-r1:7b",
    "title": "DeepSeek-r1:7b",
    "systemMessage": "You are an expert software developer. You give helpful and concise responses.",
    "apiBase": "http://localhost:11434",
    "provider": "ollama"
  },
  {
    "model": "deepseek-r1:14b",
    "title": "DeepSeek-r1:14b",
    "systemMessage": "You are an expert software developer. You give helpful and concise responses.",
    "apiBase": "http://localhost:11434",
    "provider": "ollama"
  }
],
```

```json
{
  "tabAutocompleteModel": {
    "title": "StarCoder2-3b",
    "model": "starcoder2:3b",
    "provider": "ollama"
  }
}
```

### Cline (Claude Dev)

```text
API Provider: Ollama
Model ID: llama3.2
```

### AnythingLLM

[Connecting to Ollama](https://docs.anythingllm.com/setup/llm-configuration/local/ollama)

### NextChat

[Setting Up Ollama](https://docs.nextchat.dev/models/ollama) for Seamless Integration with NextChat