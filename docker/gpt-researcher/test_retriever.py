"""
Testing your Retriever

https://docs.gptr.dev/docs/gpt-researcher/search-engines/test-your-retriever
"""

import asyncio
import os

from gpt_researcher.config.config import Config
from gpt_researcher.actions.retriever import get_retrievers
from gpt_researcher.skills.researcher import ResearchConductor
import pprint


def load_env(file_path):
    """Load environment variables from a file."""
    with open(file_path, "r", encoding="utf-8") as file:
        for line in file:
            if line.strip() and not line.startswith("#"):
                key, value = line.strip().split("=", 1)
                os.environ[key] = value


async def test_scrape_data_by_query():
    # Initialize the Config object
    config = Config()

    # Retrieve the retrievers based on the current configuration
    retrievers = get_retrievers({}, config)
    print("Retrievers:", retrievers)

    # Create a mock researcher object with necessary attributes
    class MockResearcher:
        def init(self):
            self.retrievers = retrievers
            self.cfg = config
            self.verbose = True
            self.websocket = None
            self.scraper_manager = None  # Mock or implement scraper manager
            self.vector_store = None  # Mock or implement vector store

    researcher = MockResearcher()
    research_conductor = ResearchConductor(researcher)
    print("research_conductor", dir(research_conductor))
    print("MockResearcher", dir(researcher))
    # Define a sub-query to test
    sub_query = "design patterns for autonomous ai agents"

    # Iterate through all retrievers
    for retriever_class in retrievers:
        # Instantiate the retriever with the sub-query
        retriever = retriever_class(sub_query)

        # Perform the search using the current retriever
        search_results = await asyncio.to_thread(retriever.search, max_results=10)

        print("\033[35mSearch results:\033[0m")
        pprint.pprint(search_results, indent=4, width=80)


if __name__ == "__main__":
    load_env("env_vars.txt")
    asyncio.run(test_scrape_data_by_query())
