import festival_data
import tweepy
import os
import time
import logging

### ====================================== ログ設定 ============================================
logger = logging.getLogger(__name__)

formatter = logging.Formatter(
    "%(asctime)s : %(name)s : %(levelname)s : %(lineno)s : %(message)s"
)

# 標準出力のhandlerをセット
stream_handler = logging.StreamHandler()
stream_handler.setFormatter(formatter)
stream_handler.setLevel(logging.DEBUG)

logger.addHandler(stream_handler)
# =============================================================================================


def main():
    # 環境変数郡を取得
    ## Notion関連の環境変数
    NOTION_API_TOKEN = os.getenv("NOTION_API_TOKEN")
    NOTION_DATABASE_ID = os.getenv("NOTION_DATABASE_ID")
    ## X関連の環境変数
    X_API_KEY = os.getenv("X_API_KEY")
    X_API_KEY_SECRET = os.getenv("X_API_KEY_SECRET")
    X_API_BEARER_TOKEN = os.getenv("X_API_BEARER_TOKEN")
    X_API_ACCESS_TOKEN = os.getenv("X_API_ACCESS_TOKEN")
    X_API_ACCESS_TOKEN_SECRET = os.getenv("X_API_ACCESS_TOKEN_SECRET")

    # Notionクラインアント
    notion_client = festival_data.NotionClient(api_token=NOTION_API_TOKEN)
    # X(旧twitter)クライアント
    x_client = festival_data.XClient(
        bearer_token=X_API_BEARER_TOKEN,
        consumer_key=X_API_KEY,
        consumer_secret=X_API_KEY_SECRET,
        access_token=X_API_ACCESS_TOKEN,
        access_token_secret=X_API_ACCESS_TOKEN_SECRET,
    )

    print("python batch start!!")

    # 祭り情報をXにポストする
    post_festival_data(
        notion_client=notion_client, x_client=x_client, database_id=NOTION_DATABASE_ID
    )

    time.sleep(2)

    # Xに投稿済みの祭り情報をリポストする
    quoted_repost_festival_data(
        notion_client=notion_client, x_client=x_client, database_id=NOTION_DATABASE_ID
    )

    print("python batch end!!")


def post_festival_data(notion_client, x_client, database_id):
    print("festival data post start !")
    unposted_festivals = festival_data.get_unposted(
        notion_client=notion_client, database_id=database_id
    )
    logger.debug(unposted_festivals)
    for festival in unposted_festivals:
        post_id = festival_data.post(x_client=x_client, festival=festival).get(
            "post_id"
        )
        festival_data.update_post_id(
            notion_client=notion_client, festival=festival, post_id=post_id
        )
    print("festival data post end!")


def quoted_repost_festival_data(notion_client, x_client, database_id):
    print("festival data repost start !")
    held_today_festivals = festival_data.held_today(
        notion_client=notion_client, database_id=database_id
    )
    for festival in held_today_festivals:
        repost_id = festival_data.quoted_repost(
            x_client=x_client, festival=festival
        ).get("repost_id")
        festival_data.update_repost_id(
            notion_client=notion_client, festival=festival, repost_id=repost_id
        )
    print("festival data repost end!")


if __name__ == "__main__":
    main()
