from flask import Flask, request
from gevent import pywsgi
from geventwebsocket.handler import WebSocketHandler
from flask_cors import CORS
import json

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
    return

@app.route("/api/readuser/", methods=['GET'])
def readUserHandler():
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