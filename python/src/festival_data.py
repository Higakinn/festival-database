import datetime
import requests
import textwrap
import tweepy
import pytz

from typing import List, Optional
from pydantic import BaseModel, HttpUrl


class XClient:
    """
    X(旧ツイッター)クライアントクラス
    """

    def __init__(
        self,
        bearer_token,
        consumer_key,
        consumer_secret,
        access_token,
        access_token_secret,
    ) -> None:
        self.__client: tweepy.Client = tweepy.Client(
            bearer_token=bearer_token,
            consumer_key=consumer_key,
            consumer_secret=consumer_secret,
            access_token=access_token,
            access_token_secret=access_token_secret,
        )

    def post(self, content: str):
        """
        Xにポストするための関数
        """
        _post_result = self.__client.create_tweet(text=content)
        # TODO:エラーハンドリング
        return _post_result.data.get("id")


class NotionClient:
    """
    notionクライアントクラス
    """

    def __init__(self, api_token: str) -> None:
        self.__headers = {
            "Accept": "application/json",
            "Notion-Version": "2022-06-28",
            "Content-Type": "application/json",
            "Authorization": f"{api_token}",
        }

    def query_database(self, database_id: str, db_filter: dict, limit=100):
        """
        notionデータベースをクエリするための関数
        """
        url = f"https://api.notion.com/v1/databases/{database_id}/query"
        payload = {
            "filter": db_filter,
            "page_size": limit,
        }

        response = requests.post(url, json=payload, headers=self.__headers)
        if response.status_code == requests.codes.ok:
            return response.json().get("results")

        return None

    def update_page(self, page_id: str, update_props: dict):
        """
        notionページを更新するための関数
        """
        url = f"https://api.notion.com/v1/pages/{page_id}"
        payload = {
            "properties": update_props,
        }

        response = requests.patch(url, json=payload, headers=self.__headers)
        if response.status_code == requests.codes.ok:
            return {"ok": True}

        return {"ok": False}


class Festival(BaseModel):
    """
    祭礼モデル
    """

    id: str
    name: str
    region: str
    access: str
    # start_date: date
    # end_date: Optional[date]
    start_date: str
    end_date: Optional[str]
    url: HttpUrl
    x_url: Optional[HttpUrl] = None


def get_unposted(notion_client: NotionClient, database_id) -> List[Festival]:
    """
    未投稿の祭りデータを取得するための関数
    """
    db_filter = {
        "and": [
            # NOTE: {"property": "<カラム名>", <notionプロパティ>:<該当プロパティのフィルタリング>}
            ###### 参考: https://developers.notion.com/reference/post-database-query-filter#the-filter-object
            {"property": "is_post", "checkbox": {"equals": False}}
        ]
    }
    query_db_result = notion_client.query_database(
        database_id=database_id, db_filter=db_filter
    )

    result = []
    for r in query_db_result:
        page_id = r.get("id")

        _props = r.get("properties")
        festival_name = _props.get("festival_name").get("title")[0].get("plain_text")
        region = _props.get("region").get("rich_text")[0].get("plain_text")
        access = _props.get("access").get("rich_text")[0].get("plain_text")

        _date_dict = _props.get("date").get("date")
        start_date = _date_dict.get("start")
        end_date = _date_dict.get("end")
        if _date_dict.get("end") is None:
            end_date = start_date

        url = _props.get("link").get("url")

        result.append(
            Festival(
                id=page_id,
                name=festival_name,
                region=region,
                access=access,
                start_date=start_date,
                end_date=end_date,
                url=HttpUrl(url=url),
            )
        )
    return result


def held_today(notion_client: NotionClient, database_id: str) -> List[Festival]:
    """
    実行日時に開催される祭りデータを取得するための関数
    """
    today = datetime.datetime.today(pytz.timezone("Asia/Tokyo"))
    db_filter = {
        "and": [
            # NOTE: {"property": "<カラム名>", <notionプロパティ>:<該当プロパティのフィルタリング>}
            ###### 参考: https://developers.notion.com/reference/post-database-query-filter#the-filter-object
            {"property": "is_post", "checkbox": {"equals": True}},
            {"property": "is_repost", "checkbox": {"equals": False}},
            {"property": "date", "date": {"equals": f"{today}"}},
            {"property": "date", "date": {"this_week": {}}},
        ]
    }
    query_database_result = notion_client.query_database(
        database_id=database_id, db_filter=db_filter
    )
    result = []
    for r in query_database_result:
        page_id = r.get("id")

        _props = r.get("properties")
        festival_name = _props.get("festival_name").get("title")[0].get("plain_text")
        region = _props.get("region").get("rich_text")[0].get("plain_text")
        access = _props.get("access").get("rich_text")[0].get("plain_text")
        x_url = _props.get("x url").get("formula").get("string")
        _date_dict = _props.get("date").get("date")
        start_date = _date_dict.get("start")
        end_date = _date_dict.get("end")
        url = _props.get("link").get("url")

        result.append(
            Festival(
                id=page_id,
                name=festival_name,
                region=region,
                access=access,
                start_date=start_date,
                end_date=end_date,
                url=HttpUrl(url=url),
                x_url=HttpUrl(url=x_url),
            )
        )
    return result


def _post_content(festival: Festival):
    """
    ポストする内容を生成するための関数
    """
    date = f"{festival.start_date} ~ {festival.end_date}"
    if festival.start_date == festival.end_date:
        date = festival.start_date
    post_content = (
        textwrap.dedent(
            """
【🏮祭り情報🏮】
#{festival_name}

■ 開催期間
・{date}

■ 開催場所
・{region}

■ アクセス
・{access}
■ 参考
{url}
  """
        )
        .format(
            region=festival.region,
            access=festival.access,
            festival_name=festival.name,
            date=date,
            url=festival.url,
        )
        .strip()
    )

    return post_content


def post(x_client: tweepy.Client, festival: Festival):
    """
    Xに祭り情報をポストするための関数
    """
    post_id = x_client.post(_post_content(festival))
    # TODO: エラーハンドリング
    return {"post_id": post_id}


def _quoted_repost_content(festival: Festival):
    """
    リポストする内容を生成するための関数
    """
    # TODO: 数日間にわたって開催される祭りのときに開催期間中に引用リポストできるようなポスト内容を生成できるようにする
    # today = datetime.datetime.today(pytz.timezone("Asia/Tokyo"))

    repost_content = (
        textwrap.dedent(
            """
【{region}】
#{festival_name} 始まります！

■ アクセス
・{access}

{post_url}
  """
        )
        .format(
            region=festival.region,
            access=festival.access,
            festival_name=festival.name,
            post_url=festival.x_url,
        )
        .strip()
    )
    return repost_content


def quoted_repost(x_client: XClient, festival: Festival):
    """
    Xに投稿済みの祭り情報を引用リポストするための関数
    """
    _repost_id = x_client.post(_quoted_repost_content(festival))
    # TODO: エラーハンドリング
    return {"repost_id": _repost_id}


def update_post_id(notion_client: NotionClient, festival: Festival, post_id: str):
    """
    NotionDBの該当データのpost_idカラムを更新する
    """
    update_props = {
        "is_post": {"checkbox": True},
        "post_id": {"rich_text": [{"text": {"content": post_id}}]},
    }
    return notion_client.update_page(festival.id, update_props=update_props)


def update_repost_id(notion_client: NotionClient, festival: Festival, repost_id: str):
    """
    NotionDBの該当データのrepost_idカラムを更新する
    """
    update_props = {
        "is_repost": {"checkbox": True},
        "repost_id": {"rich_text": [{"text": {"content": repost_id}}]},
    }
    return notion_client.update_page(festival.id, update_props=update_props)
