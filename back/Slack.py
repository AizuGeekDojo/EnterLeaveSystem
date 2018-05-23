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

    def postData(self, uname, sid, isent, ts, test=0, ext):
        inout = "退出" if isent else "入室"
        time = datetime.datetime.fromtimestamp(ts)
        outPutText = "{0} : {1} さんが {2} に {3} しました。{4}".format(sid, uname, time, inout, ext)
        if test:
            outPutText = "Test"
        
        self.slack.notify(text=outPutText,
                          channel=self.channel,
                          usename=self.userName,
                          icon_emoji=self.icon
                          )


if __name__ == '__main__':
    Webhook = webhookSlack(webHookURL)
    Webhook.postData("XXXXXXX", "Hoge Fugato", 0, 1526799945)
