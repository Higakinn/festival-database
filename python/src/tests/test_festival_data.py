import unittest
from pydantic import HttpUrl

import festival_data


class TestFestival(unittest.TestCase):
    def test_post_content(self):
        test_cases = [
            # まつりが１日のみの場合
            {
                "args": festival_data.Festival(
                    id="test",
                    name="開催期間が1日の祭り",
                    region="〇〇県〇〇市",
                    access="〇〇駅から〇〇分",
                    start_date="2024-07-21",
                    end_date="2024-07-21",
                    url=HttpUrl(url="http://example.com"),
                ),
                "excepted": """
【🏮祭り情報🏮】
#開催期間が1日の祭り

■ 開催期間
・2024-07-21

■ 開催場所
・〇〇県〇〇市

■ アクセス
・〇〇駅から〇〇分
■ 参考
http://example.com/
""".strip(),
            },
            {
                "args": festival_data.Festival(
                    id="test",
                    name="開催日が数日間ある祭り",
                    region="〇〇県〇〇市",
                    access="〇〇駅から〇〇分",
                    start_date="2024-07-20",
                    end_date="2024-07-22",
                    url=HttpUrl(url="http://example.com"),
                ),
                "excepted": """
【🏮祭り情報🏮】
#開催日が数日間ある祭り

■ 開催期間
・2024-07-20 ~ 2024-07-22

■ 開催場所
・〇〇県〇〇市

■ アクセス
・〇〇駅から〇〇分
■ 参考
http://example.com/
""".strip(),
            },
        ]

        for test_case in test_cases:
            args = test_case.get("args")
            excepted = test_case.get("excepted")
            exec_result = festival_data._post_content(args)
            self.assertEqual(excepted, exec_result)

    def test_quoted_repost_content(self):
        test_cases = [
            # まつりが１日のみの場合
            {
                "args": festival_data.Festival(
                    id="test",
                    name="祭礼1",
                    region="〇〇県〇〇市",
                    access="〇〇駅から〇〇分",
                    start_date="2024-07-21",
                    end_date="2024-07-21",
                    url=HttpUrl(url="http://example.com"),
                    x_url=HttpUrl(url="http://test-x_url.com"),
                ),
                "excepted": """
【〇〇県〇〇市】
#祭礼1 始まります！

■ アクセス
・〇〇駅から〇〇分

http://test-x_url.com/
""".strip(),
            },
            {
                "args": festival_data.Festival(
                    id="test",
                    name="祭礼2",
                    region="〇〇県〇〇市",
                    access="〇〇駅から〇〇分",
                    start_date="2024-07-20",
                    end_date="2024-07-22",
                    url=HttpUrl(url="http://example.com"),
                    x_url=HttpUrl(url="http://test-x_url.com"),
                ),
                "excepted": """
【〇〇県〇〇市】
#祭礼2 始まります！

■ アクセス
・〇〇駅から〇〇分

http://test-x_url.com/
""".strip(),
            },
        ]

        for test_case in test_cases:
            args = test_case.get("args")
            excepted = test_case.get("excepted")
            exec_result = festival_data._quoted_repost_content(args)
            self.assertEqual(excepted, exec_result)


if __name__ == "__main__":
    unittest.main()
