# screenshot-to-code

[screenshot-to-code](https://github.com/abi/screenshot-to-code.git)

Backend:

```bash
cd local/screenshot-to-code/backend
python3.12 -m venv .venv
source .venv/bin/activate
poetry install --no-root
pip install --upgrade fastapi pydantic
# brew install ffmpeg
```

Frontend:

```bash
cd local//screenshot-to-code/frontend
yarn
yarn dev
```
