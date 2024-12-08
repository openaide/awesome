"""
This module provides functionality to set up a connection
with OpenAI's API using environment variables for configuration.
"""

import os
import glob
from openai import OpenAI

from vanna.flask import VannaFlaskApp
from vanna.chromadb.chromadb_vector import ChromaDB_VectorStore
from vanna.openai.openai_chat import OpenAI_Chat


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

vn.connect_to_postgres(
    host=os.getenv("POSTGRES_HOST", "localhost"),
    port=os.getenv("POSTGRES_PORT", "5432"),
    dbname=os.getenv("POSTGRES_DBNAME", "postgres"),
    user=os.getenv("POSTGRES_USER", "postgres"),
    password=os.getenv("POSTGRES_PASSWORD", ""),
)

#
vn.train(sql="SELECT version();")
vn.train(sql="SELECT * FROM pg_catalog.pg_user;")
vn.train(sql="SELECT * FROM pg_database WHERE datistemplate = false AND datallowconn = true;")
vn.train(sql="SELECT * FROM information_schema.tables;")

# read the sql from TRAIN_PATH and train the model
train_path = os.getenv("TRAIN_PATH", "/workspace/train")
sql_files = glob.glob(os.path.join(train_path, '*.sql'))
for sql_file in sql_files:
    with open(sql_file, 'r', encoding='utf-8') as file:
        sql_content = file.read()
        vn.train(sql=sql_content)

app = VannaFlaskApp(vn)
host = os.environ.get("HOST", "0.0.0.0")
port = int(os.environ.get("PORT", 5000))
app.run(host=host, port=port)
