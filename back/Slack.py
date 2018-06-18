import os
import settings
import slackweb
import datetime


webHookURL = os.getenv("SLACK_WEBHOOK_URL")

class webhookSlack():
    def __init__(self, URL, userName="Logging Bot", icon=":robot_face:", channel="#enter_leave_log"):
        self.slack = slackweb.Slack(url=URL)
        self.userName = userName
        self.icon = icon
        self.channel = channel

    def postData(self, uname, sid, isent, ts, ext, test=0):
        inout = "退出" if isent else "入室"
        time = datetime.datetime.fromtimestamp(ts)
        outPutText = "{0} : {1} さんが {2} に {3} しました。".format(sid, uname, time, inout)

        self.postAgdIO(outPutText)

    def postAgdIO(self, msg):
        self.slack.notify(text=msg,
                          channel=self.channel,
                          usename=self.userName,
                          icon_emoji=self.icon
                          )

    def postAgdOut(self, msg, ext):
        purpose = ext['Use']
        impress = ext['message']

        attachments = []
        elements = {"title": "アンケート結果",
                    "pretext": msg,
                    "text": "目的 : {0} \n 感想 : {1}".format(purpose, impress)
                    }

        attachments.append(elements)
        self.slack.notify(attachments=attachments,
                          channel=self.channel,
                          usename=self.userName,
                          icon_emoji=self.icon
                          )


if __name__ == '__main__':
    Webhook = webhookSlack(webHookURL)
    jsonForOut = {'Use': ['Other'], 'message': 'This is Test'}
    jsonForIn = {}
    Webhook.postData("XXXXXXX", "Hoge Fugato", 0, 1526799945, jsonForIn)
