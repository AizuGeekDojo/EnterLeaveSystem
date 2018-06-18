import os
import time
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

        if isent == 0:
            self.postAgdIn(outPutText)

        elif isent:
            self.postAgdOut(outPutText, ext)

    def postAgdIn(self, msg):
        self.slack.notify(text=msg,
                          channel=self.channel,
                          usename=self.userName,
                          icon_emoji=self.icon
                          )

    def postAgdOut(self, msg, ext):
        purpose = ", ".join([e for e in ext['Use']])
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
    jsonForOut1 = {'Use': ['3DPrinter', 'LaserCutter'], 'message': '3DPrinter Test'}
    jsonForOut2 = {'Use': ['LaserCutter', 'Other'], 'message': 'LaserCutter Test'}
    jsonForOut3 = {'Use': ['Training session', '3DPrinter'], 'message': 'Training session Test'}
    jsonForIn = {}
    Webhook.postData("XXXXXXX", "Hoge Fugato", 1, 1526799945, jsonForOut1)
    time.sleep(3)
    Webhook.postData("XXXXXXX", "Hoge Yasaka", 1, 1526799945, jsonForOut2)
    time.sleep(3)
    Webhook.postData("XXXXXXX", "Noah Orberg", 1, 1526799945, jsonForOut3)
    time.sleep(3)

