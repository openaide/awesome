# syntax=docker/dockerfile:1

ARG AIDER_BASE_IMAGE

###
# https://aider.chat/docs/faq.html#how-can-i-run-aider-locally-from-source-code
# https://aider.chat/docs/llms/openai-compat.html
FROM ${AIDER_BASE_IMAGE} AS aider-gui

COPY ./local/aider /aider

WORKDIR /aider
RUN python -m pip install -e .

ENV PIP_NO_INPUT=1
RUN python -m pip install --upgrade --upgrade-strategy only-if-needed 'aider-chat[browser]'

WORKDIR /workspace
ENTRYPOINT []
CMD [ "python", "-m", "aider", "--gui" ]
###
