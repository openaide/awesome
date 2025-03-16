#!/bin/bash

uv venv
source .venv/bin/activate 

uv pip install -r requirements.txt
uv pip install -r ./local/OpenManus/requirements.txt
uv pip install -e ./local/OpenManus