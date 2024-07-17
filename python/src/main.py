import festival_data
import tweepy
import os


def main():
    print("python batch start!!")
    NOTION_API_TOKEN = os.getenv("NOTION_API_TOKEN")
    NOTION_DATABASE_ID = os.getenv("NOTION_DATABASE_ID")
    festivals = festival_data.get_all(NOTION_API_TOKEN, NOTION_DATABASE_ID)

    X_API_KEY = os.getenv("X_API_KEY")
    X_API_KEY_SECRET = os.getenv("X_API_KEY_SECRET")
    X_API_BEARER_TOKEN = os.getenv("X_API_BEARER_TOKEN")
    X_API_ACCESS_TOKEN = os.getenv("X_API_ACCESS_TOKEN")
    X_API_ACCESS_TOKEN_SECRET = os.getenv("X_API_ACCESS_TOKEN_SECRET")

    # 認証
    client = tweepy.Client(
        bearer_token=X_API_BEARER_TOKEN,
        consumer_key=X_API_KEY,
        consumer_secret=X_API_KEY_SECRET,
        access_token=X_API_ACCESS_TOKEN,
        access_token_secret=X_API_ACCESS_TOKEN_SECRET,
    )
    for festival in festivals:
        region = festival.get("region")
        access = festival.get("access")
        festival_name = festival.get("festival_name")
        date = festival.get("date")
        url = festival.get("url")
        page_id = festival.get("page_id")
        print(festival_name)
        post_id = festival_data.post_data(
            client=client,
            region=region,
            access=access,
            festival_name=festival_name,
            date=date,
            url=url,
        ).get("post_id")
        print(post_id)
        festival_data.update(NOTION_API_TOKEN, page_id, post_id)

    print("python batch end!!")


if __name__ == "__main__":
    main()
