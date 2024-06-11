from grpc import insecure_channel
from json import dumps
from django.http import HttpResponse
from django.views.decorators.csrf import csrf_exempt
from django.conf import settings

from proto.cargo_pb2_grpc import CargoHoldStub
from proto.cargo_pb2 import ContentReply, ContentRequest
from stoke.django_annotations import require_token
from stoke.grpc_client import GRPCClientCredentials

def index(request):
    if request.user.is_anonymous:
        return HttpResponse(f"Hello unknown user!".encode())
    return HttpResponse(f"Hello {request.user}".encode())

@csrf_exempt
@require_token(claims={"inv" : "acc"})
def test(request):
    return HttpResponse(b'{ "hello": "world", "foo" : "bar" }', status=200, content_type="application/json")

@csrf_exempt
@require_token(claims={"inv" : "acc"})
def cargo_contents(request):
    creds = GRPCClientCredentials(token=request.META["STOKE_AUTH_TOKEN"])
    # TODO Must be over secure channel: https://github.com/grpc/grpc/issues/33618
    with insecure_channel(settings.CARGO_GRPC_ADDRESS) as channel:
        stub = CargoHoldStub(channel) 
        reply : ContentReply = stub.GetContents(ContentRequest(), credentials=creds.call_credentials())
        return HttpResponse(('{ "contents" : ' + dumps([item.name for item in reply.items]) + '}').encode(), status=200, content_type="application/json")


