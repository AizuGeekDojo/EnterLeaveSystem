from flask import Flask, request
from gevent import pywsgi
from geventwebsocket.handler import WebSocketHandler
from flask_cors import CORS
import json
from handler import *

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
        print('input responce ')
        msg = json.dumps({"message": True})
        ws.send(msg)
    return 

@app.route("/api/createuser/", methods=['POST'])
def createUserHandler():
    req_json = json.loads(request.data.decode('utf-8'))
    createUser(req_json)
    return

@app.route("/api/readuser/", methods=['GET'])
def readUserHandler():
    req_json = json.loads(request.data.decode('utf-8'))
    getUser(req_json)
    return

@app.route("/api/updateuser", methods=['UPDATE'])
def updateUserHandler():
    return

def main():
    app.debug = True
    server = pywsgi.WSGIServer(("", 3000), app, handler_class=WebSocketHandler)
    print("server runnning at port:3000")
    server.serve_forever()

if __name__ == "__main__":
	main()