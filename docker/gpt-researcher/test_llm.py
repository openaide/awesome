"""
Testing your LLM

https://docs.gptr.dev/docs/gpt-researcher/llms/testing-your-llm
"""

import asyncio
import os

from gpt_researcher.config.config import Config
from gpt_researcher.utils.llm import create_chat_completion


def load_env(file_path):
    """Load environment variables from a file."""
    with open(file_path, "r", encoding="utf-8") as file:
        for line in file:
            if line.strip() and not line.startswith("#"):
                key, value = line.strip().split("=", 1)
                os.environ[key] = value


async def main():
    cfg = Config()

    try:
        report = await create_chat_completion(
            model=cfg.smart_llm_model,
            messages=[{"role": "user", "content": "sup?"}],
            temperature=0.35,
            llm_provider=cfg.smart_llm_provider,
            stream=True,
            # max_tokens=cfg.max_tokens,
            llm_kwargs=cfg.llm_kwargs,
        )
        print(report)
    except Exception as e:
        print(f"Error in calling LLM: {e}")


if __name__ == "__main__":
    load_env("env_vars.txt")
    asyncio.run(main())
