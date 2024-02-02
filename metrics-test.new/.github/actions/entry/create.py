import github
import os
from datetime import datetime, timedelta
from dotenv import load_dotenv
import pandas as pd
import requests
import json
from notion_client import Client

def create_db_entry():
    notion_token = os.getenv('NOTION_KEY')
    db_id = os.getenv('NOTION_DATABASE_ID')
    client = Client(auth=notion_token)
    title = os.getenv('TITLE')

    new_page = client.pages.create(
        **{
            "parent": { "database_id": db_id },
            "properties": {
                "名前": {
                    "title": [
                        {
                            "text": {
                                "content": title
                            }
                        }
                    ]
                }
            }
        }
    )
    new_database = client.databases.create(
        **{
            "parent": { "page_id": new_page['id'] },
            "is_inline": True,
            "title": [
                {
                    "text": {
                        "content": "index"
                    }
                }
            ],
            "properties": {
                "Name": {
                    "title": {}
                },
                "Date": {
                    "type": 'created_time',
                    "created_time": {}
                },
                "Tags": {
                    "multi_select": {}
                }
            }
        }
    )
    return new_database['id']

def main():
    load_dotenv()
    database_id = '\ndatabase_id=' + create_db_entry() + '\n'
    output = os.getenv('GITHUB_OUTPUT')
    with open(output, mode='a') as f:
        f.write(database_id)

if __name__ == "__main__":
    main()
