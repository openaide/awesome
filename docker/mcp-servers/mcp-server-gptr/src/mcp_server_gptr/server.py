"""
MCP server for generating research reports
"""

import json
from enum import Enum
from dotenv import load_dotenv

from mcp.server.models import InitializationOptions
from mcp.types import ImageContent, EmbeddedResource, TextContent, Tool
from mcp.server import NotificationOptions, Server
from mcp.server.stdio import stdio_server

from pydantic import BaseModel
from gpt_researcher import GPTResearcher
from gpt_researcher.utils.enum import ReportType, Tone


from .backend.detailed_report import DetailedReport


class GptrTools(str, Enum):
    """
    GPT Researcher Tools
    """

    GENERATE_REPORT = "generate_report"


class GptrResult(BaseModel):
    """
    The result of the research report generation
    """

    report: str
    costs: float


class GptrInput(BaseModel):
    """
    The input parameters for the research report generation
    Attributes:
        query (str): The user query to generate the report.
        report_type (str): The type of report to generate.
        tone (str): The tone of the report.
    """

    query: str
    report_type: str
    tone: str


async def generate_report(args: GptrInput) -> GptrResult:
    """
    Generate a research report based on the input parameters.
    """
    tone_map = {
        "objective": Tone.Objective,
        "formal": Tone.Formal,
        "analytical": Tone.Analytical,
        "persuasive": Tone.Persuasive,
        "informative": Tone.Informative,
        "explanatory": Tone.Explanatory,
        "descriptive": Tone.Descriptive,
        "critical": Tone.Critical,
        "comparative": Tone.Comparative,
        "speculative": Tone.Speculative,
        "reflective": Tone.Reflective,
        "narrative": Tone.Narrative,
        "humorous": Tone.Humorous,
        "optimistic": Tone.Optimistic,
        "pessimistic": Tone.Pessimistic,
    }
    tone = tone_map[args.tone] if args.tone in tone_map else Tone.Objective
    type_map = {
        "summary": ReportType.ResearchReport,
        "detailed": ReportType.DetailedReport,
        "deep": ReportType.DeepResearch,
    }
    report_type = (
        type_map[args.report_type]
        if args.report_type in type_map
        else ReportType.ResearchReport
    )
    if args.report_type == "detailed":
        detailed_report = DetailedReport(
            query=args.query,
            report_type="research_report",
            report_source="web_search",
            source_urls=[],
            tone=tone,
            subtopics=[],
            headers={},
        )
        detailed_report.gpt_researcher.set_verbose(True)

        report = await detailed_report.run()

        costs = detailed_report.gpt_researcher.get_costs()
    else:
        researcher = GPTResearcher(
            query=args.query,
            report_type=report_type,
            tone=tone,
        )
        researcher.set_verbose(True)

        await researcher.conduct_research()
        report = await researcher.write_report()

        costs = researcher.get_costs()
    return GptrResult(report=report, costs=costs)


server = Server("mcp-server-gptr")


@server.list_tools()
async def handle_list_tools() -> list[Tool]:
    """
    List available research report tools.
    """
    return [
        Tool(
            name=GptrTools.GENERATE_REPORT.value,
            description="""Generate factual and impartial research reports,
complete with citations, derived from extensive web online research""",
            inputSchema={
                "type": "object",
                "properties": {
                    "query": {
                        "type": "string",
                        "description": "The user query to generate the report",
                    },
                    "report_type": {
                        "type": "string",
                        "description": """The type of report to generate. Options include:
summary: Short and fast report
detailed: In depth and longer well-structured report
deep: Comprehensive and thorough with depth and breadth""",
                    },
                    "tone": {
                        "type": "string",
                        "description": """The tone of the report. Defaults to 'objective'. Options include:
objective: Impartial and unbiased presentation
formal: Academic standards with sophisticated language
analytical: Critical evaluation and examination
persuasive: Convincing viewpoint
informative: Clear and comprehensive information
explanatory: Clarifying complex concepts
descriptive: Detailed depiction
critical: Judging validity and relevance
comparative: Juxtaposing different theories
speculative: Exploring hypotheses
reflective: Personal insights
narrative: Story-based presentation
humorous: Light-hearted and engaging
optimistic: Highlighting positive aspects
pessimistic: Focusing on challenges""",
                    },
                },
                "required": ["query", "report_type"],
            },
        )
    ]


@server.call_tool()
async def handle_call_tool(
    name: str, arguments: dict | None
) -> list[TextContent | ImageContent | EmbeddedResource]:
    """
    Handle tool calls for research report request.
    """
    try:
        match name:
            case GptrTools.GENERATE_REPORT.value:
                query = arguments.get("query")
                report_type = arguments.get("report_type")
                tone = arguments.get("tone")

                if not query or not report_type:
                    raise ValueError("Missing query or report_type")

                input_params = GptrInput(
                    query=query, report_type=report_type, tone=tone
                )
                result = await generate_report(input_params)
            case _:
                raise ValueError(f"Unknown tool: {name}")

        await server.request_context.session.send_resource_list_changed()

        return [
            TextContent(
                type="text",
                text=result.report,
                metadata={
                    "costs": result.costs,
                    "report_type": report_type,
                    "tone": tone,
                },
            )
        ]

    except Exception as e:
        raise ValueError(f"Error processing query: {str(e)}") from e


async def main():
    """
    Main function to start the server.
    """
    async with stdio_server() as (read_stream, write_stream):
        load_dotenv()
        await server.run(
            read_stream,
            write_stream,
            InitializationOptions(
                server_name="gptr",
                server_version="0.1.0",
                capabilities=server.get_capabilities(
                    notification_options=NotificationOptions(),
                    experimental_capabilities={},
                ),
            ),
        )
