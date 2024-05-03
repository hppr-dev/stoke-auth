from flask import Flask
from stoke.flask_annotations import require_claims
from stoke.client import StokeClient

app = Flask(__name__)
client = StokeClient("http://localhost:8080/api/pkeys")

@app.route("/fire")
@require_claims(client, { "role" : "eng" }, jwt_kwarg="token")
def hello_world(token=dict()):
    return f"<p>pew pew! authorized by { token['n'] }!</p>"
