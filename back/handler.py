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
    sid = req_json["SID"]
    card_id = req_json["CardID"]
    user_name = db.getUserName(sid)
    res = json.dumps({
        "SID": sid,
        "CardID": card_id,
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


def isNewCard(card_id: str) -> bool:
    """
    新しいカードかの確認
    """
    return db.getSIDByIDm(card_id) is None

