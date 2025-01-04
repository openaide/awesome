"""
https://github.com/huggingface/smolagents
"""

from smolagents import CodeAgent, DuckDuckGoSearchTool, HfApiModel, LiteLLMModel

model = LiteLLMModel(model_id="gpt-4o")

agent = CodeAgent(tools=[DuckDuckGoSearchTool()], model=model)

agent.run(
    "How many seconds would it take for a leopard at full speed to run through Pont des Arts?"
)
