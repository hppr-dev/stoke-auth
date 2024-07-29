# Stoke Python Client Library

# Stoke Client

A StokeClient object is used to pull signing certificates from the stoke server.

```python
# No ssl
client = StokeClient("http://172.17.0.1:8080")

# configuring ssl
ssl_context = SSLContext(PROTOCOL_TLS_CLIENT)
ssl_context.load_verify_locations("/app/examples/certs/ca.crt")
client = StokeClient("https://172.17.0.1:8080", ssl_context=ssl_context)
```

# Flask require_claims annotation

The @require_claims annotation requires tokens with the specified claims.
Set the jwt_kwarg to allow passing the parsed token as a kwarg as the receiving function.

```python
from stoke.flask_annotations import require_claims

@app.route("/shipment")
@require_claims(client, { "req" : "acc" }, jwt_kwarg="token")
def request_item(token=dict()):
    return "Shipment requested"
```

# Django

## Auth Backend
The django auth backend allows users access to the django admin page and using stoke tokens to authorize access.

The backend requires the following settings:
  * STOKE_AUTH_CONFIG  -- a dictionary defining the configuration for the stoke server.

The STOKE_AUTH_CONFIG dictionary MUST have the following keys:
  * CLIENT : StokeClient | TestStokeClient -- the client to use when verifying tokens
  * BASE_CLAIMS : dict[str, str] -- the minimum claims that a user needs to authenticate
  * USERNAME : str -- the claim key to use as a user's username

Optional settings:
  * ROUTE_CLAIMS : dict[str, dict[str, str]] -- a dictionary of route paths and their required claims
  * STAFF: 2-tuple -- the claim (key, value) to use to determine if a user is staff
  * SUPERUSER : 2-tuple -- the claim (key, value) to use to determine if a user is a superuser
  * FIRST_NAME : str -- the claim key to lookup a user's first name
  * LAST_NAME: str -- the claim key to lookup a user's last name
  * EMAIL: str -- the claim key to lookup a user's email

```python
stoke_client =  StokeClient("http://172.17.0.1:8080")
STOKE_AUTH_CONFIG = {
    'CLIENT' : stoke_client,
    'BASE_CLAIMS': { "inv" : "acc" },
    'USERNAME' : 'u',
    'EMAIL': 'e',
    'SUPERUSER' : ('stk', 'S'),
    'STAFF' : ('stk', 'S'),
}
```

## require_claims Annotation

```python 
from stoke.django_annotations import require_claims

@require_token(claims={"inv" : "acc"})
def endpoint(request):
    return HttpResponse(b'Hello world')
```

# GRPC

## Server

Use intercept_all to require claims for all GRPC server endpoints.

```python
from stoke.grpc_server import intercept_all
stoke_client =  StokeClient("http://172.17.0.1:8080")
server = new_server(
    thread_pool=futures.ThreadPoolExecutor(max_workers=10),
    interceptors=intercept_all(stoke_client, required_claims={"car":"acc"})
)
```

## Client

```python 
from stoke.grpc_client import GRPCClientCredentials

my_token = "SOMETOKEN"
call_creds = GRPCClientCredentials(token=my_token).call_credentials()

with secure_channel(settings.CARGO_GRPC_ADDRESS, chan_creds) as channel:
    stub = MyGRPCStub(channel) 
    reply = stub.CallGRPCFunction(GRPCRequest(), credentials=call_creds)
```
