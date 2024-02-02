# 開発チームの活動状況の概要をGitHubより得る

## 事前準備

1. Notionのインテグレーションの作成と適用
    - 開発者向けには「インテグレーション」、通常のNotion上では「コネクト」と名称が違うので注意
    - https://notion.so/my-integrations で新しいインテグレーションを作成する
        - 一意な名称を設定する
        - ここで生成される「シークレット」が本リポジトリに設定する`NOTION_KEY`
        - 「機能」をクリック、チェックボックスすべて（5つ）をチェックする
    - 新しいデータベースを作成
        - 「新規ページ」をクリック
        - 新しく開いたページでまず「テーブル」を選択
        - ページの名称を設定する
        - 先にページの名和尚を入れてEnterしてしまうと、ページになってしまうので注意
    - コネクト
        - ウィンドウの右上の「・・・」をクリックし下方にある「コネクトの追加」を選択し、表示される一覧から、作成したインテグレーションの名称を選択
        - このページのURLのDBのID部分をリポジトリの`NOTION_DATABASE_ID`
            - 例）URLが `https://www.notion.so/infomart/53d19556cb9246afb9eb099417aaa139?v=aae1811e9b054166b36e1e61f3c931bd`である場合、`https://www.notion.so/infomart/`の次から`?`までの間の`53d19556cb9246afb9eb099417aaa139`がID

2. GitHub Actionsのためのシークレットを設定（本リポジトリをもとに新規にリポジトリを作成する場合）

    - MY_ACCESS_TOKEN
        - データを収集したいリポジトリにアクセス可能なパーソナルアクセストークン
            - b2bplatformに参加している開発者のPATを作成すること
    - NOTION_KEY
        - Notionのインテグレーションのシークレット
    - NOTION_DATABASE_ID
        - Notionへの書き込みの親となるデータベースのID
    
    ※ 抽出したデータを表示するためのインデックスとなるNotionのデータベースIDは、シークレット化していない

## 特徴

あくまでもサマリーとしてチーム（実際はリポジトリ単位）の活動状況を提示する（細かな情報は、Projectsで取ってください）
- GitHubのサマライズ
    - 提示する情報は以下の通り
        - Issues / pr / Discussions
            - 最初のコメント、またはレビューがなされるまでの時間
            - クローズまでの時間
            - 解答までの時間（Discussionsのみ）
            - ラベルが削除されるまでの時間（指定時のみ）
            - 「時間」は、平均値・中央値、90パーセンタイルを提示
            - 上記とは別に、オープンのままの数、クローズされた数、作成された数を提示
            - 現時点では、月次・週次で以下を抽出
                - 作成されたIssue
                - 作成されたpr
                - クローズされたIssue
                - クローズされたpr
        - Commit
            - コミット数を人別日別に提示
- すべての情報はNotionのデータベースまたはページにて提示
    - Metircs Index（仮称）のDBにリポジトリごとの月次・週次のインデックスを列挙
    - 上記インデックスからリポジトリごとのインデックス
    - リポジトリごとのインデックスから各情報
- GitHub Actionsにて月次・週次で起動
    - GitHubは内部ではUTCで時刻管理しているため、ローカルタイムへの変換が必要
        - 現時点では、毎月初日の2:03（日本時間11:03）、および毎日曜の2:03（日本時間11:03）に起動
            - キリの良い時間に設定しないの原則だが、もう少し考えるべきか？
            - 前日の15:00が日本時間の0時あるから、週次なら土曜に設定で対処できるが、月次は月ごと、うるう年などの考慮が必要なため（つまりは面倒であるため）上記の時間で起動をかけている

## 用意したプライベートアクション

### commit

- 与えられた期間のリポジトリのコミット情報を抽出し、日別人別のコミット数として整形し、Notion DBに登録する
- v3 APIを使用
- Pythonで記述、コンテナで動作させる
- 入力（環境変数）
    - ACCESS_TOKEN
        - GitHubのパーソナルアクセストークン（`secrets.MY_ACCESS_TOKEN`）
        - 複数のリポジトリにアクセスするのなら`secrets.GITHUB_TOKEN`ではなく個人のプライベートアクセストークンにする
    - NOTION_KEY
        - Notionのインテグレーションシークレット（`secrets.NOTION_KEY`）
    - NOTION_DATABASE_ID
        - 親となるDBのID（`secrets.NOTION_DATABASE_ID`）
    - FROM_DATE
        - 情報抽出期間の初日（UTC）
    - TO_DATE
        - 情報抽出期間の最終日（UTC）
    - REPO_OWNER
        - リポジトリにオーナ/組織
    - REPO_NAME
        - リポジトリ名
    - TITLE
        - 本情報のタイトル
- 出力
    なし
- 備考
    - UTCはISO 8601拡張形式とする（YYYY-MM-DDThh:mm:ssZ）
        - 内部で、与えられたTimeZoneの初日から最終日で処理している（つまり時刻情報は捨てている）

### cutter

- 与えられたファイルの指定した空白行の間のテキストを削除する
- Pythonで記述、コンテナで動作

- 入力（環境変数）
    - INPUT_FILE
        - 処理対象のファイル
    - OUTPUT_FILE
        - 処理後の出力ファイル（INPUT_FILEと同じ場合、上書き）
    - START_BANK_LINE
        - 削除を開始する空白行の序数（何番目の空白行）
    - END_BLANK_LINE
        - 削除を終了する空白行の序数（何番目の空白行）
- 出力
    なし

### entry

- Notionでインラインでデータベースを持つページを作成する
- Pythonで記述、コンテナで動作

- 入力（環境変数）
    - NOTION_KEY
        - Notionのインテグレーションシークレット（`secrets.NOTION_KEY`）
    - NOTION_DATABASE_ID
        - 親となるDBのID（`secrets.NOTION_DATABASE_ID`）
    - REPO_OWNER
        - リポジトリのオーナ/組織
    - REPO_NAME
        - リポジトリ名
    - TITLE
        - 本ページとDBのタイトル
    
### localdate

- 与えたUTC日時をもとに、与えたTimeZoneに対応した前月、または前週の始まりと終わりをUTCで出力
- Pythonで記述、コンテナで動作させる

- 入力（環境変数）
    - UTC:
        - UTC日時
    - TYPE
        - 'month'または'week'
    - WEEKDAY
        - 週の始まりの曜日（デフォルトは日曜）
    - TIMEZONE
        - 変換するタイムゾーン
- 出力
    - ${{ steps.<実行したステップのid>.outputs.first }}
        月、または週の最初の日時（UTC）
    - ${{ steps.<実行したステップのid>.outputs.last }}
        月、または週の最後の日時（UTC）
- 備考
    - UTCはISO 8601拡張形式とする（YYYY-MM-DDThh:mm:ssZ）

### notion

- Notionのページを作成し、MarkdownファイルをNotionのブロック形式に変換し提示
    - github/issue-metrics@v2はMarkdownおよびJSONで結果を出力するが、Notion APIにはファイルのアップロードやインポート機能が存在しないため作成
    - Notionの制約でブロック内の項目が100を超えるとエラーとなる（メッセージ表示のみでGitHub Actionsの実行は継続）
- JavaScriptで記述、そのまま実行

- 入力（環境変数）
    - NOTION_KEY
        - Notionのインテグレーションシークレット（`secrets.NOTION_KEY`）
    - NOTION_DATABASE_ID
        - 親となるDBのID（`secrets.NOTION_DATABASE_ID`）
    - MARKDOWN_FILENAME
        - Markdownファイルのパス
    - TITLE
        - 本情報のタイトル
    - TAGS
        - インデックスとなるDBにつけるタグ
- 備考
    - 動作形態を考える必要があるかもしれない
        - 「そのまま実行」であるため、node_modulesをアップロードしているのが気持ち悪い
        - コンテナで動作させるか（処理速度は落ちそう）

## 使用するアクション

### [github/issue-metrics@v2](https://github.com/github/issue-metrics)

- 指定したリポジトリのIssues/Pull Requests/Discussionsの統計情報をサマライズする（詳細は「特徴」の項参照）
- 結果は`issue_metrics.md`または`issue_metrics.json`でファイル出力
- GitHub製でPythonで記述、コンテナで実行（内部ではまさかのv3 APIを使用）

- 入力（環境変数）（使用しているもののみ列挙）
    - GH_TOKEN
        - リポジトリにアクセスするためのトークン（`secres.MY_ACCESS_TOKEN`を使用）
    - SEARCH_QUERY
        - Issues/Pull Requests/Discussionsから抽出する条件で、以下を使用
            - `repo:b2bplatform/<リポジトリ名>`
            - `is:issue`または`is:pr`
            - `created:<日時>..<日時>`または`closed:<日時>..<日時>`
- 出力
    なし
- 備考
    - 日時はISO 8601拡張形式でUTCで指定（YYYY-MM-DDThh:mm:ssZ）
