import time
import db
import json

def createUser(req_json: dict):
    """
    ユーザーの作成
    """

    sid = req_json["SID"]
    card_id = req_json["CardID"]

    db.addUser(card_id,sid)
    success = True

    res = json.dumps({
        "Success": success,
        "SID": sid,
        "CardID": card_id,
        "timestamp": int(time.time())
    })
    return res

def getUser(req_json: dict):
    """
    ユーザーの取得
    """
    card_id = req_json["CardID"]
    sid = db.getSIDByIDm(card_id)
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

    db.updateUser(card_id,sid)
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
    ext = str(req_json["Ext"])
    ts = req_json["timestamp"]

    db.addLog(sid,isent,ext,ts)

    res = json.dumps({
        "SID": sid,
        "timestamp": int(time.time())
    })
    return res

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
