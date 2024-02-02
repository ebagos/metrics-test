import github
import os
from datetime import datetime, timedelta
from dotenv import load_dotenv
import pandas as pd
import requests
import json
from notion_client import Client
import zoneinfo
from dateutil.tz import gettz

# --------------------------------------------
def github_commit_data():
    # 環境変数の取得
    tzone = gettz(os.getenv('TIMEZONE'))

    start_date = datetime.strptime(os.getenv("FROM_DATE"), '%Y-%m-%dT%H:%M:%SZ')
    start_date = start_date.astimezone(tz=tzone)
    
    end_date = datetime.strptime(os.getenv("TO_DATE"), '%Y-%m-%dT%H:%M:%SZ')
    end_date = end_date.astimezone(tz=tzone)

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
        # 現時点でページング処理を行っていないため、100件以上のコミットがある場合は取得できない
        branch_commits = repo.get_commits(since=start_date, until=end_date, sha=branch.commit.sha, )

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
                # authorやcommitterがない場合がある
                if commit.author == None:
                    commit_data['author'].append(commit.sha)
                    # 日付だけのデータにする
                    day = datetime.strptime(end_date.strftime('%Y-%m-%d'), '%Y-%m-%d')
                    commit_data['date'].append(day)
                else:
                    commit_data['author'].append(commit.author.login)
                    # 日付だけのデータにする
                    day = datetime.strptime(commit.commit.author.date.strftime('%Y-%m-%d'), '%Y-%m-%d')
                    commit_data['date'].append(day)
                commit_data['message'].append(commit.commit.message)
                commit_data['url'].append(commit.html_url)
                commit_data['branch'].append(branch)

    # --------------------------------------------
    # Create a Pandas DataFrame from the dictionary
    df = pd.DataFrame.from_dict(commit_data)
    
    if df.size == 0:
        data_count = pd.DataFrame({'author': [0]})
    else:
        # 時刻を丸めて日付でグループ化し、authorごとの数を集計
        data_count = df.groupby([df['date'].dt.date, 'author']).size().unstack(fill_value=0)

    # 前日付についてデータフレームを再構築
    sd = start_date.strftime('%Y-%m-%d')
    ed = end_date.strftime('%Y-%m-%d')
    data_count = data_count.reindex(pd.date_range(start=sd, end=ed), fill_value=0)

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
    send_df_to_notion(df)

if __name__ == "__main__":
    main()
