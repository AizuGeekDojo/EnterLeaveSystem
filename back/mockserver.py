from flask import Flask, request
from gevent import pywsgi
from geventwebsocket.handler import WebSocketHandler
import json
app = Flask(__name__)

@app.route("/")
def index():
    return "Hello World!"

@app.route('/socket')
def socket():
    if request.environ.get('wsgi.websocket'):
        print("Connected")
        ws = request.environ['wsgi.websocket']
        while True:
            msg = json.dumps({"message": input()})
            ws.send(msg)

def main():
    app.debug = True
    server = pywsgi.WSGIServer(("", 3000), app, handler_class=WebSocketHandler)
    print("server runnning at port:3000")
    server.serve_forever()
    
if __name__ == "__main__":
	main()


