import asyncio
import time
import os

from app.agent.manus import Manus
from app.logger import logger

os.environ["ANONYMIZED_TELEMETRY"] = "false"


async def main():
    agent = Manus()
    try:
        prompt = input("Enter your prompt: ")
        if not prompt.strip():
            logger.warning("Empty prompt provided.")
            return

        logger.warning("Processing your request...")
        start_time = time.time()
        await agent.run(prompt)
        logger.info("Request processing completed.")
        end_time = time.time()
        elapsed_time = end_time - start_time
        logger.info(f"Elapsed time: {elapsed_time:.2f} seconds.")
    except KeyboardInterrupt:
        logger.warning("Operation interrupted.")


if __name__ == "__main__":
    asyncio.run(main())
