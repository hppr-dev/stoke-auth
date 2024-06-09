from typing import Any, Dict
from json import loads, dumps
from urllib.request import HTTPError, Request, urlopen
from datetime import timedelta, datetime
from base64 import urlsafe_b64decode

def request_token(url, username, password) -> Dict[Any, Any]:
    req = Request(
        url,
        dumps({ "username" : username, "password": password }).encode(),
        { "Content-Type" : "application/json" },
        method="POST"
    )
    try:
        with urlopen(req) as response:
            return loads(response.read())
    except HTTPError:
        return {"token": "", "refresh": ""}

def refresh_token(url, token, refresh) -> Dict[Any, Any]:
    req = Request(
        url,
        dumps({ "refresh": refresh }).encode(),
        { "Content-Type" : "application/json", "Authorization" : f"Bearer {token}" },
        method="POST"
    )
    try:
        with urlopen(req) as response:
            return loads(response.read())
    except HTTPError:
        return {"token": "", "refresh": ""}

def expires_from_token(token: str) -> datetime:
    claim_str = token.split(".")[1]
    claims = loads(urlsafe_b64decode(claim_str + ( '=' * ( 4 - len(claim_str) % 4 ) ) ))
    if "exp" in claims:
        return datetime.fromtimestamp(claims["exp"])
    return datetime.now() + timedelta(minutes=5)
