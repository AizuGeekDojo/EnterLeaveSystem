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

@app.route('/socket')
def socket():
    if request.environ.get('wsgi.websocket'):
        print("Connected")
        ws = request.environ['wsgi.websocket']
        print('input responce ')
        msg = json.dumps({"message": True})
        ws.send(msg)

    return "hoge"

@app.route('/api/getUser')
def getUser():
	res = json.dumps({
		"isExist": True,
		"StudentId": "s1240000",
		"Name": "Hoge Fuga"
	})
	return res

# sample request curl  -H "Content-type: application/json" -X POST localhost:3000/api/register -d '{"hoge": "fuga"}'
@app.route('/api/register', methods=['POST'])
def register():
    req_json = json.loads(request.data.decode('utf-8'))
    print(req_json)
    return json.dumps({'error': False, 'data': req_json})	

def main():
    app.debug = True
    server = pywsgi.WSGIServer(("", 3000), app, handler_class=WebSocketHandler)
    print("server runnning at port:3000")
    server.serve_forever()
    
if __name__ == "__main__":
	main()


