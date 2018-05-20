from flask import Flask, request
from gevent import pywsgi
from geventwebsocket.handler import WebSocketHandler
from flask_cors import CORS
import json
from handler import *
import nfc_read

app = Flask(__name__)
CORS(app)

@app.route("/")
def index():
    return "Hello World!"

@app.route('/socket/readCard')
def socket():
    if request.environ.get('wsgi.websocket'):
        print("Connected")
        ws = request.environ['wsgi.websocket']
        while True:
            cardid = nfc_read.nfc_read()
            print("||| " + cardid + " ")
            if cardid is not "":
                msg = json.dumps({
                    "IsCard": True,
                    "CardID": cardid,
                    "IsNew": isNewCard(cardid),
                    "timestamp": int(time.time())
                })
                ws.send(msg)
                break
    return 

@app.route("/api/createuser", methods=['POST'])
def createUserHandler():
    req_json = json.loads(request.data.decode('utf-8'))
    res = createUser(req_json)
    return res

@app.route("/api/readuser", methods=['POST'])
def readUserHandler():
    req_json = json.loads(request.data.decode('utf-8'))
    res = getUser(req_json)
    return res

@app.route("/api/updateuser", methods=['UPDATE'])
def updateUserHandler():
    return

class WebSocket(self):
    def open_websocket(self):
        app.debug = True
        self.server = pywsgi.WSGIServer(("", 3000), app, handler_class=WebSocketHandler)
        print("server runnning at port:3000")
        server.serve_forever()
    def close_websocket(self):
        self.server.close()


if __name__ == "__main__":
    ws = WebSocket()
	ws.open_websocket()
