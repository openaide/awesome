"""
This app sets up a local context for Vanna's OpenAI model and ChromaDB
"""

import argparse
import logging
import os
import sys

from openai import OpenAI
from vanna.chromadb.chromadb_vector import ChromaDB_VectorStore
from vanna.flask import VannaFlaskApp
from vanna.openai.openai_chat import OpenAI_Chat

log_level = os.getenv("LOG_LEVEL", "INFO").upper()
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
                "path": os.getenv("STORE_PATH", "./local/store"),
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
        dbname=os.getenv("POSTGRES_DBNAME", ""),
        user=os.getenv("POSTGRES_USER", ""),
        password=os.getenv("POSTGRES_PASSWORD", ""),
    )


# https://github.com/vanna-ai/vanna/blob/main/src/vanna/base/base.py#L1865
# https://github.com/vanna-ai/notebooks/blob/main/postgres-openai-vanna-vannadb.ipynb
def basic_training():
    """
    Train the model with basic SQL commands.
    """
    vn.train(sql="SELECT version();")
    vn.train(sql="SELECT * FROM pg_catalog.pg_user;")
    vn.train(sql="SELECT * FROM pg_catalog.pg_database;")

    # df_information_schema = vn.run_sql("""
    #     SELECT * FROM INFORMATION_SCHEMA.COLUMNS
    # """)
    # plan = vn.get_training_plan_generic(df_information_schema)
    # vn.train(plan=plan)


# vn.train(ddl=ddl, documentation=doc, sql=sql)
def find_and_train(train_path):
    """
    Find and train the model with the provided training data.
    """
    file_types = ["ddl", "documentation", "sql"]

    for file_type in file_types:
        file_type_path = os.path.join(train_path, file_type)
        if not os.path.exists(file_type_path):
            # logging.warning("Directory does not exist: %s", file_type_path)
            continue

        for root, _, filenames in os.walk(file_type_path):
            for filename in filenames:
                file_path = os.path.join(root, filename)
                logging.debug("Training on file: %s", file_path)
                try:
                    with open(file_path, "r", encoding="utf-8") as file:
                        content = file.read()
                        vn.train(**{file_type: content})
                        logging.debug("Training successful on file: %s", file_path)
                except (IOError, OSError) as e:
                    logging.error("Error reading file: %s. Error: %s", file_path, e)


def is_training_done():
    """
    Check if the training is already done.
    """
    store_path = os.getenv("STORE_PATH", "./local/store")
    train_done = os.path.join(store_path, ".train_done")
    return os.path.exists(train_done)


def set_training_done():
    """
    Set the training as done.
    """
    store_path = os.getenv("STORE_PATH", "./local/store")
    if not os.path.exists(store_path):
        os.makedirs(store_path)
    train_done = os.path.join(store_path, ".train_done")
    with open(train_done, "w", encoding="utf-8") as file:
        file.write("")


def train_model():
    """
    Train the model with basic SQL commands and the provided training data.
    """
    train_path = os.getenv("TRAIN_PATH", "./local/train")
    basic_training()
    find_and_train(train_path)
    set_training_done()


class CustomArgumentParser(argparse.ArgumentParser):
    """Custom argument parser for handling command-line argument parsing
    with personalized error messages in the Vanna Application.
    """

    def error(self, message):
        sys.stderr.write(f"Error: {message}\n")
        sys.stderr.write("You must specify one of train or serve\n")
        self.print_help()
        sys.exit(2)


def main():
    """Set up and parse command line arguments for the Vanna Application."""
    parser = CustomArgumentParser(description="Vanna Application")

    parser.add_argument("--dbname", type=str, help="Postgres database name")
    parser.add_argument("--user", type=str, help="Postgres user")
    parser.add_argument("--password", type=str, help="Postgres password")

    subparsers = parser.add_subparsers(
        dest="command", required=True, help="Available commands"
    )

    # Sub-parser for the "train" command
    subparsers.add_parser("train", help="Start training")

    # Sub-parser for the "serve" command
    serve_parser = subparsers.add_parser("serve", help="Start the server")
    serve_parser.add_argument(
        "--host",
        type=str,
        help="Specify the host address",
        default=os.environ.get("HOST", "0.0.0.0"),
    )
    serve_parser.add_argument(
        "--port",
        type=int,
        help="Specify the port number",
        default=os.environ.get("PORT", 5000),
    )

    #
    args = parser.parse_args()

    if not args.command:
        parser.print_help(sys.stderr)
        sys.exit(1)

    if args.dbname:
        os.environ["POSTGRES_DBNAME"] = args.dbname
    if args.user:
        os.environ["POSTGRES_USER"] = args.user
    if args.password:
        os.environ["POSTGRES_PASSWORD"] = args.password
    if not (
        os.getenv("POSTGRES_DBNAME")
        and os.getenv("POSTGRES_USER")
        and os.getenv("POSTGRES_PASSWORD")
    ):
        print(
            "Usage: python app.py --dbname <dbname> --user <user> "
            "--password <password> ...\n\n"
            "Missing required arguments: dbname, user, or password\n\n"
            "Arguments for database can be provided via command line "
            " or set as environment variables:\n"
            "POSTGRES_DBNAME, POSTGRES_USER, POSTGRES_PASSWORD.\n\n"
            "All arguments are required either from command line"
            " or environment variables.\n"
            "You may also set POSTGRES_HOST and POSTGRES_PORT."
            " default: localhost:5432"
        )
        sys.exit(1)

    if args.command == "train":
        init_db()
        train_model()
        logging.info("Training completed successfully")
    elif args.command == "serve":
        init_db()
        if not is_training_done():
            train_model()
        app = VannaFlaskApp(vn, title="Welcome", allow_llm_to_see_data=True, debug=True)
        logging.info("Listening on %s:%s", args.host, args.port)
        app.run(host=args.host, port=args.port)


if __name__ == "__main__":
    main()
