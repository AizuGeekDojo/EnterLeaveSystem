import time
import db


def createUser(req_json: dict):
    """
    ユーザーの作成
    """
    #TODO

    sid = req_json["SID"]
    card_id = req_json["CardID"]

    addUser(card_id,sid)
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
    #TODO
    user_name = pass
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
    #TODO
    
    sid = req_json["SID"]
    card_id = req_json["CardID"]
    
    updateUser(card_id,sid)
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
    return getSIDByIDm(card_id)

