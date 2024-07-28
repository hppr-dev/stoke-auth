from os import environ
from ssl import PROTOCOL_TLS_CLIENT, SSLContext
from flask import Flask

from stoke.flask_annotations import require_claims
from stoke.client import StokeClient
from stoke.test_client import TestStokeClient

app = Flask(__name__)

if "STOKE_TEST" in environ:
    client = TestStokeClient(default_dict={"req": "acc"})
else:
    ssl_context = SSLContext(PROTOCOL_TLS_CLIENT)
    ssl_context.load_verify_locations("/app/examples/certs/ca.crt")
    client = StokeClient("https://172.17.0.1:8080", ssl_context=ssl_context)

@app.route("/shipment")
@require_claims(client, { "req" : "acc" }, jwt_kwarg="token")
def request_item(token=dict()):
    return "Shipment requested"
