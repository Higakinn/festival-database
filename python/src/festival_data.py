import requests
import textwrap


def get_all(api_token, database_id, limit=100):
    # notion get
    url = f"https://api.notion.com/v1/databases/{database_id}/query"
    # is_post=false つまり  まだ投稿していない祭礼情報を取得するためのフィルタ
    db_filter = {
        "and": [
            # {
            #     "property": "region"
            # },
            {"property": "is_post", "checkbox": {"equals": False}}
        ]
    }
    payload = {
        "filter": db_filter,
        "page_size": limit,
    }
    print(payload)
    headers = {
        "Accept": "application/json",
        "Notion-Version": "2022-06-28",
        "Content-Type": "application/json",
        "Authorization": f"{api_token}",
    }

    response = requests.post(url, json=payload, headers=headers).json()
    print(response)
    result = []
    for r in response.get("results"):
        page_id = r.get("id")

        _props = r.get("properties")
        festival_name = _props.get("festival_name").get("title")[0].get("plain_text")
        region = _props.get("region").get("rich_text")[0].get("plain_text")
        access = _props.get("access").get("rich_text")[0].get("plain_text")

        _date_dict = _props.get("date").get("date")
        _start_date = _date_dict.get("start")
        _end_date = _date_dict.get("end")
        date = f"{_start_date} ~ {_end_date}"
        if _end_date is None:
            date = _start_date

        url = _props.get("link").get("url")

        result.append(
            {
                "page_id": page_id,
                "festival_name": festival_name,
                "region": region,
                "access": access,
                "date": date,
                "url": url,
            }
        )
    return result


def post_data(client, region, access,festival_name, date, url):
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
            region=region,
            access=access,
            festival_name=festival_name,
            date=date,
            url=url,
        )
        .strip()
    )
    _post_result = client.create_tweet(text=post_content)
    return {"post_id": _post_result.data.get("id")}


def update(api_token, page_id, post_id):
    url = f"https://api.notion.com/v1/pages/{page_id}"
    update_props = {
        "is_post": {"checkbox": True},
        "post_id": {"rich_text": [{"text": {"content": post_id}}]},
    }
    payload = {
        "properties": update_props,
    }
    headers = {
        "Accept": "application/json",
        "Notion-Version": "2022-06-28",
        "Content-Type": "application/json",
        "Authorization": f"{api_token}",
    }

    response = requests.patch(url, json=payload, headers=headers)
    if response.status_code == requests.codes.ok:
        return {"ok": True}

    return {"ok": False}
