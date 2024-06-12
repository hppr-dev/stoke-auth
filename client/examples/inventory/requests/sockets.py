from queue import SimpleQueue
from urllib.parse import parse_qs
from django.conf import settings
from channels.generic.websocket import JsonWebsocketConsumer
from grpc import secure_channel, ssl_channel_credentials

from proto.cargo_pb2 import Item, LoadItemRequest
from proto.cargo_pb2_grpc import CargoHoldStub
from stoke.grpc_client import GRPCClientCredentials

class LoadItemConsumer(JsonWebsocketConsumer):
    def connect(self):
        qs=parse_qs(self.scope['query_string'].decode())
        if "token" not in qs or len(qs['token']) != 1:
            return
        call_creds = GRPCClientCredentials(token=qs['token'][0]).call_credentials()
        chan_creds = ssl_channel_credentials(
            root_certificates=settings.CA_FILE_DATA,
            private_key=settings.KEY_FILE_DATA,
            certificate_chain=settings.CERT_FILE_DATA,
        )
        # Must be over secure channel: https://github.com/grpc/grpc/issues/33618
        self.channel = secure_channel(settings.CARGO_GRPC_ADDRESS, chan_creds)
        stub = CargoHoldStub(self.channel) 
        self.send_q = SimpleQueue()
        self.receiver = stub.LoadItems(iter(self.send_q.get, None), credentials=call_creds)
        self.accept()

    def receive_json(self, content, **kwargs):
        item = Item(name=content["name"], id=content["id"], height=1, width=1, depth=1)
        self.send_q.put(LoadItemRequest(num_items=content["num"], item=item))
        recv = next(self.receiver)
        self.send_json({"loaded": recv.loaded, "message" : recv.message})

    def disconnect(self, close_code):
        self.channel.close()
        self.receiver.cancel()


