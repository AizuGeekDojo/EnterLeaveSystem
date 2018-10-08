import time
import db
import json
import settings
import os
import slackweb
import Slack

webHookURL = os.getenv("SLACK_WEBHOOK_URL")


def createUser(req_json: dict):
    """
    ユーザーの作成
    """

    sid = req_json["SID"]
    card_id = req_json["CardID"]

    if db.getUserName(sid) is None:
        success = False
    else:
        db.addUser(card_id,sid)
        success = True


    res = json.dumps({
        "Success": success,
        "SID": sid,
        "CardID": card_id,
        "timestamp": int(time.time())
    })
    return res


def getUser(sid: dict):
    """
    ユーザーの取得
    """
    user_name = db.getUserName(sid)
    is_enter = isEnter(sid)
    res = json.dumps({
        "SID": sid,
        "CardID": card_id,
        "IsEnter": is_enter,
        "UserName": user_name,
        "timestamp": int(time.time())
    })
    return res


def updateUser(req_json: dict):
    """
    ユーザの更新
    """

    sid = req_json["SID"]
    card_id = req_json["CardID"]

    db.updateUser(card_id, sid)
    success = True

    res = json.dumps({
        "SID": sid,
        "CardID": card_id,
        "timestamp": int(time.time())
    })
    return res


def addLog(req_json: dict):

    sid = req_json["SID"]
    isent = req_json["IsEnter"]
    ext = req_json["Ext"]
    ts = req_json["timestamp"]
    uname = db.getUserName(sid)

    db.addLog(sid, isent, ts, json.dumps(ext))

    res = json.dumps({"SID": sid, "timestamp": int(time.time())})
    return res


def slack_notify(req_json: dict):
    sid = req_json["SID"]
    isent = req_json["IsEnter"]
    ext = req_json["Ext"]
    ts = req_json["timestamp"]
    uname = db.getUserName(sid)
    slack = Slack.webhookSlack(URL=webHookURL)
    slack.postData(uname, sid, isent, int(ts / 1000), ext)


def isNewCard(card_id: str) -> bool:
    """
    新しいカードかの確認
    """
    return db.getSIDByIDm(card_id) is None


def isEnter(sid: str) -> bool:
    """
    入室済みかの確認
    """
    return db.isUserInside(sid)


def getSID(card_id: str):
    """
    カードIDからsidに変換
    """
    return db.getSIDByIDm(card_id)