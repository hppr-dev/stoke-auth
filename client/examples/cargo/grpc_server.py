from os import environ
from concurrent import futures
from ssl import PROTOCOL_TLS_CLIENT, SSLContext
from grpc import server as new_server, ssl_server_credentials
from grpc_reflection.v1alpha import reflection
from stoke.client import StokeClient
from stoke.grpc_server import intercept_all
from proto.cargo_pb2 import ContentReply, ContentRequest, Item, LoadItemReply, LoadItemRequest, DESCRIPTOR as CARGO_DESCRIPTOR
from proto.cargo_pb2_grpc import CargoHoldServicer, add_CargoHoldServicer_to_server

class CargoHoldServer(CargoHoldServicer):

    def __init__(self, x, y, z):
        self.dim = (x, y , z)
        self.contents = []

    def GetContents(self, request: ContentRequest, context):
        print("Sending contents...")
        return ContentReply(
            filter_matched=True,
            items=self.contents,
        )

    def LoadItems(self, request_iterator, context):
        print("Loading items...")
        for req in request_iterator:
            yield self._load(req)
            if not context.is_active:
                break
        print("Done loading items.")

    def _load(self, req : LoadItemRequest) -> LoadItemReply:
        for i in range(req.num_items):
            self.contents.append(req.item)
        return LoadItemReply(loaded=True)

def read_file(filename : str) -> bytes :
    with open(filename, 'rb') as f:
        return f.read()

cert_file = environ.get("CARGO_CERT")
key_file  = environ.get("CARGO_KEY")
ca_file   = environ.get("CARGO_CA")

ssl_context = SSLContext(PROTOCOL_TLS_CLIENT)
ssl_context.load_verify_locations(ca_file)

stoke_client = StokeClient(
    url="https://172.17.0.1:8080",
    ssl_context=ssl_context,
)

def serve():
    print("Starting grpc server...")
    server = new_server(
        thread_pool=futures.ThreadPoolExecutor(max_workers=10),
        interceptors=intercept_all(stoke_client, required_claims={"car":"acc"})
    )

    add_CargoHoldServicer_to_server(CargoHoldServer(100, 100, 100), server)

    SERVICE_NAMES = (
        CARGO_DESCRIPTOR.services_by_name['CargoHold'].full_name,
        reflection.SERVICE_NAME,
    )
    reflection.enable_server_reflection(SERVICE_NAMES, server)

    if cert_file is not None and key_file is not None and ca_file is not None:
        print("Reading certificate files...")
        credentials = ssl_server_credentials(
            [(read_file(key_file), read_file(cert_file))],
            root_certificates=read_file(ca_file),
            require_client_auth=True,
        )
        server.add_secure_port("0.0.0.0:6060", credentials)
    else:
        print("Could not read CARGO_CERT, CARGO_KEY, or CARGO_CA. Starting in insecure mode (python client credentials will not work).")
        server.add_insecure_port("0.0.0.0:6060")


    server.start()
    print("Listening for connections...")
    server.wait_for_termination()

if __name__ == "__main__":
    serve()
