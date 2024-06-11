from concurrent import futures
from grpc import server as new_server
from grpc_reflection.v1alpha import reflection
from stoke.client import StokeClient
from stoke.grpc_server import intercept_all
from proto.cargo_pb2 import ContentReply, ContentRequest, LoadItemReply, LoadItemRequest, DESCRIPTOR as CARGO_DESCRIPTOR
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
        print("Done loading items.")

    def _load(self, req : LoadItemRequest) -> LoadItemReply:
        self.contents.append(req)
        return LoadItemReply(loaded=True)


def serve():
    print("Starting grpc server...")
    server = new_server(
        thread_pool=futures.ThreadPoolExecutor(max_workers=10),
        interceptors=intercept_all(StokeClient(url="http://172.17.0.1:8080"), required_claims={"car":"acc"})
    )
    add_CargoHoldServicer_to_server(CargoHoldServer(100, 100, 100), server)
    SERVICE_NAMES = (
        CARGO_DESCRIPTOR.services_by_name['CargoHold'].full_name,
        reflection.SERVICE_NAME,
    )
    reflection.enable_server_reflection(SERVICE_NAMES, server)
    server.add_insecure_port("0.0.0.0:6060")
    server.start()
    print("Listening for connections...")
    server.wait_for_termination()

if __name__ == "__main__":
    serve()
