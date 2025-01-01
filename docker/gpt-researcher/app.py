from gpt_researcher import GPTResearcher

import asyncio
import os
import json
from uuid import uuid4


def load_env_vars(file_path):
    """Load environment variables from a file."""
    with open(file_path, "r", encoding="utf-8") as file:
        for line in file:
            if line.strip() and not line.startswith("#"):
                key, value = line.strip().split("=", 1)
                os.environ[key] = value


async def generate_report(query: str, report_type: str):
    """
    Generate a research report based on the provided query and report type.
    """
    researcher = GPTResearcher(query, report_type)
    researcher.set_verbose(True)
    await researcher.conduct_research()
    report = await researcher.write_report()

    # Get additional information
    research_context = researcher.get_research_context()
    source_urls = researcher.get_source_urls()
    costs = researcher.get_costs()
    research_images = researcher.get_research_images()
    research_sources = researcher.get_research_sources()

    return {
        "report": report,
        "costs": costs,
        "research_context": research_context,
        "source_urls": source_urls,
        "research_images": research_images,
        "research_sources": research_sources,
    }


def save_report(report, base_dir=None):
    """
    Write the report to a file
    """
    base_dir = base_dir or os.getcwd()
    artifact_filepath = os.path.join(base_dir, f"outputs/{uuid4()}.md")
    os.makedirs("outputs", exist_ok=True)
    with open(artifact_filepath, "w", encoding="utf-8") as f:
        f.write(report)
    print(f"Report written to '{artifact_filepath}'")


def main():
    """
    Executes an async function to run a research report.
    """
    # query = "what team may win the NBA finals?"
    query = "Help me plan an adventure to California"
    # query = "What happened in the latest burning man floods?"

    # "research_report": "Summary - Short and fast (~2 min)"
    # "detailed_report": "Detailed - In depth and longer (~5 min)"
    report_type = "research_report"

    #
    load_env_vars("env_vars.txt")
    result = asyncio.run(generate_report(query, report_type))

    print(json.dumps(result, indent=4))

    if result["report"]:
        save_report(result["report"])
    else:
        print("No report generated.")

    # print("Report:")
    # print(result.report)
    # print("\nResearch Costs:")
    # print(result.costs)
    # print("\nResearch Context:")
    # print(result.research_context)
    # print("\nSource URLs:")
    # print(result.source_urls)
    # print("\nNumber of Research Images:")
    # print(len(result.research_sources))
    # print("\nNumber of Research Sources:")
    # print(len(result.research_sources))
    # print("\n")


if __name__ == "__main__":
    main()
