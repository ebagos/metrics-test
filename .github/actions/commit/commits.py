import github
import os
from datetime import datetime, timedelta
from dotenv import load_dotenv
import pandas as pd
import requests
import json
from notion_client import Client

# --------------------------------------------
def github_commit_data():
    # 環境変数の取得
    start_date = datetime.strptime(os.getenv("FROM_DATE"), '%Y-%m-%dT%H:%M:%SZ')
    end_date = datetime.strptime(os.getenv("TO_DATE"), '%Y-%m-%dT%H:%M:%SZ')
    token = os.getenv('ACCESS_TOKEN')
    owner = os.getenv('REPO_OWNER')
    repo_name = os.getenv('REPO_NAME')

    # Authenticate to GitHub using a personal access token
    g = github.Github(token)

    # Specify the repository and time period
    repository_name = owner + '/' + repo_name

    # Get the repository
    repo = g.get_repo(repository_name)

    # Create a dictionary to store commits
    commits = {}

    # Iterate through all branches
    for branch in repo.get_branches():
        branch_name = branch.name

        # Get commits for the branch within the specified time period
        branch_commits = repo.get_commits(since=start_date, until=end_date, sha=branch.commit.sha)

        # Store the commits in the dictionary
        commits[branch_name] = branch_commits

    g.close()
    # --------------------------------------------

    # Create a dictionary to store the commit data
    commit_data = {
        'sha': [],
        'author': [],
        'date': [],
        'message': [],
        'url': [],
        'branch': []
    }

    duplicates = {}
    # Iterate through the commits dictionary
    for branch, branch_commits in commits.items():
        for commit in branch_commits:
            if commit.sha not in duplicates:
                duplicates[commit.sha] = True
                # Store the commit data in the dictionary
                commit_data['sha'].append(commit.sha)
                commit_data['author'].append(commit.author.login)
                commit_data['date'].append(commit.commit.author.date)
                commit_data['message'].append(commit.commit.message)
                commit_data['url'].append(commit.html_url)
                commit_data['branch'].append(branch)

    # --------------------------------------------
    # Create a Pandas DataFrame from the dictionary
    df = pd.DataFrame.from_dict(commit_data)

    # Convert the date column to a datetime object
    df['date'] = pd.to_datetime(df['date'])

    # 日付でグループ化し、authorごとの数を集計
    data_count = df.groupby([df['date'].dt.date, 'author']).size().unstack(fill_value=0)

    # 全日付についてデータフレームを再構築
    data_count = data_count.reindex(pd.date_range(start=start_date, end=end_date), fill_value=0)

    return data_count

# --------------------------------------------
def send_df_to_notion(df):
    notion_token = os.getenv('NOTION_KEY')
    db_id = os.getenv('NOTION_DATABASE_ID')
    client = Client(auth=notion_token)
    title = os.getenv('TITLE')
    json = df.to_dict(orient='records')

    dict = {}
    dict['date']={"title":{}}
    for item in list(df.columns.values):
        dict[item]={"number":{}}
 
    new_page = client.pages.create(
        **{
            "parent": { "database_id": db_id },
            "properties": {
                "Name": {
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
                        "content": title
                    }
                }
            ],
            "properties": dict
        }
    )
    for data in df.iterrows():
        day = data[0].strftime('%Y-%m-%d')
        dict = {}
        dict['date']={"title":[{"text":{"content":day}}]}
        for item in list(df.columns.values):
            dict[item]={"number":int(data[1][item])}
        new_page = client.pages.create(
            **{
                "parent": { "database_id": new_database['id'] },
                "properties": dict
            }
        )
    return new_database

def main():
    load_dotenv()
    df = github_commit_data()
    print(df)
    send_df_to_notion(df)

if __name__ == "__main__":
    main()
