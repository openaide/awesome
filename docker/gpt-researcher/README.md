# GPT Researcher

[GPT Researcher](https://github.com/assafelovic/gpt-researcher) is an autonomous agent designed for comprehensive web and local research on any given task.

By default, local [SearXNG](https://github.com/searxng/searxng.git) [AGPL-3.0 license](https://github.com/searxng/searxng?tab=AGPL-3.0-1-ov-file#readme) is used.
You can switch to other [search engines](https://github.com/assafelovic/gpt-researcher/blob/master/docs/docs/gpt-researcher/search-engines/retrievers.md)

See [configuration](https://github.com/assafelovic/gpt-researcher/blob/master/docs/docs/gpt-researcher/gptr/config.md) for more info.

## Run with Docker

```bash
make clone
make build
make up
make down
```

After starting up the services, visit [http://gptr.localhost/](http://gptr.localhost/)

Reports are saved in in `./outputs`. You may change it to a different location in [compose.overrid.yaml](./compose.override.yml)

## Run with CLI

```bash
make setup
make cli
# docker compose -f compose.searxng.yml up -d
# python cli.py "What are the main causes of climate change?" --report_type research_report
# python cli.py "The impact of artificial intelligence on job markets" --report_type detailed_report
# python cli.py "Renewable energy sources and their potential" --report_type outline_report
```

## Standalone agent app

```bash
make setup
make app
```
