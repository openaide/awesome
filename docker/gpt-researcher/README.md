# GPT Researcher

[GPT Researcher](https://github.com/assafelovic/gpt-researcher) is an autonomous agent designed for comprehensive web and local research on any given task.

By default, local [SearXNG](https://github.com/searxng/searxng.git) [AGPL-3.0 license](https://github.com/searxng/searxng?tab=AGPL-3.0-1-ov-file#readme) is used.
You can switch to other [search engines](https://github.com/assafelovic/gpt-researcher/blob/master/docs/docs/gpt-researcher/search-engines/retrievers.md)

See [configuration](https://github.com/assafelovic/gpt-researcher/blob/master/docs/docs/gpt-researcher/gptr/config.md) for more info.

```bash
make clone
make build
make up
make down
```

After starting up the services, visit [http://gptr.localhost/](http://gptr.localhost/)
