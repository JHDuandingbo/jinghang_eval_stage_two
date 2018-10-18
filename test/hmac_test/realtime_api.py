import time
import json
import threading
import uuid
from urllib.parse import quote
from hashlib import sha1
import hmac
import base64
from websocket import create_connection  # pip install websocket-client
import websocket

appId = "TiD3p6"
accessKeyId = "HGTBv4hFj9"
accessKeySecret = "JZ5J39vFncv3j3453X2G45sCy6cOv5G3"
base_uri = 'wss://api.iflyrec.com/ast?lang=en&codec=pcm_s16le&bitrate=16000&authString='
end_tag = "{\"end\": true}"
file_path = './1_.pcm'


def get_local_datetime():
    return time.strftime("%Y-%m-%dT%H:%M:%S+0800", time.localtime(time.time()))


def generate_signature(auth_string):
    s = quote(auth_string, encoding='utf-8')
    print("urlencoded basestring:<" + s+">")
    print("Hmac in:<" + s +">")
    my_sign_bytes = hmac.new(accessKeySecret.encode('utf-8'), s.encode('utf-8'), sha1).digest()
    print(str(my_sign_bytes))
    my_sign = str(base64.b64encode(my_sign_bytes), 'utf-8')
    print("Hmac out:<" + my_sign +">")
    return my_sign


def get_uri():
    auth_string = "v1.0," + appId + "," + accessKeyId + "," + str(get_local_datetime()) + "," + str(uuid.uuid4())
    print("basestring:<" + auth_string+">")
    signature = generate_signature(auth_string)
    final_auth_string = quote(auth_string + ',' + signature, encoding='utf-8')
    uri = base_uri + final_auth_string
    return uri


class Client(object):
    def __init__(self):
        # 得到uri
        self.uri = get_uri()
        print("URI:" + self.uri)
        self.ws = create_connection(self.uri)
        print(self.ws.connected)
        # self.trecv = threading.Thread(target=self.recv)
        # self.trecv.start()

    def send(self, f_path):
        file_object = open(f_path, 'rb')
        try:
            index = 1
            while True:
                chunk = file_object.read(1280)
                # for b in chunk:
                #     print(hex(b), end=';;')
                #     print(b, end=' ')
                # print()
                # print('---------------------------')
                # print(chunk)
                if not chunk:
                    break
                self.ws.send(chunk)

                index += 1
                time.sleep(0.04)
        finally:
            print(str(index) + ", read len:" + str(len(chunk)) + ", file tell:" + str(file_object.tell()))
            file_object.close()

        self.ws.send(bytes(end_tag, encoding='utf-8'))
        print("send end tag success")

    def recv(self):
        try:
            while self.ws.connected:
                result = str(self.ws.recv())
                if len(result) == 0:
                    print("receive result end")
                    break
                result_dict = json.loads(result)

                # 解析结果
                if result_dict["action"] == "started":
                    print("handshake success, result: " + result)

                if result_dict["action"] == "result":
                    print("rtasr result: " + result)

                if result_dict["action"] == "error":
                    print("rtasr error: " + result)
                    self.ws.close()
                    return

                print(result)
        except websocket.WebSocketConnectionClosedException:
            print("receive result end exception")

    def close(self):
        self.ws.close()
        print("connection closed")


if __name__ == '__main__':
    client = Client()
    client.send(file_path)
