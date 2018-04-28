from flask import Flask, request
from gevent import pywsgi
from geventwebsocket.handler import WebSocketHandler

app = Flask(__name__)

@app.route("/")
def index():
    return "Hello World!"

@app.route('/socket')
def socket():
    if request.environ.get('wsgi.websocket'):
        ws = request.environ['wsgi.websocket']
        while True:
            ws.send(input())

def main():
    app.debug = True
    server = pywsgi.WSGIServer(("", 3000), app, handler_class=WebSocketHandler)
    server.serve_forever()

if __name__ == "__main__":
	main()


