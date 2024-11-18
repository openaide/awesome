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
docker exec ollama ollama run llama3.2
docker exec ollama ollama run llama3.2:1b
docker exec ollama ollama run starcoder2:3b
docker exec ollama ollama run codellama:7b-instruct
docker exec ollama ollama run codellama:7b-code
docker exec ollama ollama run qwen2.5-coder:7b-base
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
      "model": "llama3.2",
      "title": "Ollama3.2",
      "systemMessage": "You are an expert software developer. You give helpful and concise responses.",
      "apiKey": "x",
      "apiBase": "http://localhost:11434",
      "provider": "ollama"
    }
]
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