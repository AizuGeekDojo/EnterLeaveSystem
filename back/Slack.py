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

    def postData(self, sid, isent, ts, test=0):
        inout = "入室" if isent else "退出"
        time = datetime.datetime.fromtimestamp(ts)
        outPutText = "{0} さんが {1} に {2} しました。".format(sid, time, inout)
        if test:
            outPutText = "Test"
        
        self.slack.notify(text=outPutText,
                          channel=self.channel,
                          usename=self.userName,
                          icon_emoji=self.icon
                          )

        return outPutText


if __name__ == '__main__':
    Webhook = webhookSlack(webHookURL)
    Webhook.postData("Hoge Fugato", 0, 1526799945)
