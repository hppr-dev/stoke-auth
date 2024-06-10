from os import environ
from flask import Flask

from stoke.flask_annotations import require_claims
from stoke.client import StokeClient
from stoke.test_client import TestStokeClient

app = Flask(__name__)

if "STOKE_TEST" in environ:
    client = TestStokeClient(default_dict={"req": "acc"})
else:
    client = StokeClient("http://172.17.0.1:8080")

@app.route("/shipment")
@require_claims(client, { "req" : "acc" }, jwt_kwarg="token")
def request_item(token=dict()):
    return "Shipment requested"
