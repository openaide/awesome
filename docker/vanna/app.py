"""
This module provides functionality to set up a connection
with OpenAI's API using environment variables for configuration.
"""

import argparse
import logging
import os

from openai import OpenAI
from vanna.chromadb.chromadb_vector import ChromaDB_VectorStore
from vanna.flask import VannaFlaskApp
from vanna.openai.openai_chat import OpenAI_Chat

log_level = os.getenv('LOG_LEVEL', 'INFO').upper()
log = logging.getLogger()
log.setLevel(log_level)


class LocalContextOpenAI(ChromaDB_VectorStore, OpenAI_Chat):
    """
    LocalContextOpenAI is a combined implementation of ChromaDB_VectorStore
    and OpenAI_Chat. This class initializes both parent classes using the
    provided configuration settings.

    :param config: An optional dictionary containing configuration parameters
        for the base classes.
    """

    def __init__(self, config=None):
        ChromaDB_VectorStore.__init__(
            self,
            config={
                "client": "persistent",
                "path": os.getenv("STORE_PATH", "./local"),
            },
        )
        client = OpenAI(
            api_key=os.getenv("LLM_API_KEY", "sk-1234"),
            base_url=os.getenv("LLM_BASE_URL", "http://localhost:4000"),
        )
        model_name = os.getenv("LLM_MODEL", "gpt-4o")
        config = {"model": model_name}
        OpenAI_Chat.__init__(self, client=client, config=config)

    def search_tables_metadata(
        self, engine, catalog, schema, table_name, ddl, size, config=None, **kwargs
    ):
        raise NotImplementedError


vn = LocalContextOpenAI()


def init_db():
    """
    Initialize the database connection using environment variables.
    """
    vn.connect_to_postgres(
        host=os.getenv("POSTGRES_HOST", "localhost"),
        port=os.getenv("POSTGRES_PORT", "5432"),
        dbname=os.getenv("POSTGRES_DBNAME", "postgres"),
        user=os.getenv("POSTGRES_USER", "postgres"),
        password=os.getenv("POSTGRES_PASSWORD", ""),
    )

    # at least one training is required to workaround a flask UI bug
    vn.train(sql="SELECT version();")


# https://github.com/vanna-ai/vanna/blob/main/src/vanna/base/base.py#L1865
# https://github.com/vanna-ai/notebooks/blob/main/postgres-openai-vanna-vannadb.ipynb
def basic_training():
    """
    Train the model with basic SQL commands.
    """
    vn.train()

    df_information_schema = vn.run_sql("""
        SELECT * FROM INFORMATION_SCHEMA.COLUMNS
    """)
    plan = vn.get_training_plan_generic(df_information_schema)
    vn.train(plan=plan)


# vn.train(ddl=ddl, documentation=doc, sql=sql)
def find_and_train(train_path):
    """
    Find and train the model with the provided training data.
    """
    file_types = ["ddl", "doc", "sql"]

    for file_type in file_types:
        for root, _, filenames in os.walk(os.path.join(train_path, file_type)):
            for filename in filenames:
                file_path = os.path.join(root, filename)
                logging.debug("training on file: %s", file_path)
                with open(file_path, "r", encoding="utf-8") as file:
                    content = file.read()
                    vn.train(**{file_type: content})


def run_training():
    """
    Run the training process.
    """
    train_path = os.getenv("TRAIN_PATH", "./local/train")
    find_and_train(train_path)


def train_model():
    """
    Train the model with basic SQL commands and the provided training data.
    """
    basic_training()
    run_training()


if __name__ == "__main__":
    # Set up argument parser
    parser = argparse.ArgumentParser(description="Run the Vanna Flask app.")
    parser.add_argument(
        "--host",
        type=str,
        help="Specify the host address",
        default=os.environ.get("HOST", "0.0.0.0"),
    )
    parser.add_argument(
        "--port",
        type=int,
        help="Specify the port number",
        default=os.environ.get("PORT", 5000),
    )
    parser.add_argument(
        "--skip-training", action="store_true", help="Skip training process"
    )

    # Parse command-line arguments
    args = parser.parse_args()

    # Initialize the app
    init_db()
    app = VannaFlaskApp(vn)

    # Handle skip training logic
    if not args.skip_training:
        train_model()

    #
    logging.info("Listening on %s:%s", args.host, args.port)
    app.run(host=args.host, port=args.port)
