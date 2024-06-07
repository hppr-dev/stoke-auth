from base64 import urlsafe_b64decode
from json import dumps, loads
from datetime import datetime, timedelta
from urllib.request import Request, urlopen
from grpc import access_token_call_credentials

class GRPCClientCredentials:
    def __init__(self, token : str, refresh_token : str = "", refresh_url : str = "", refresh_window : timedelta = timedelta(seconds=5)):
        self.token = token
        self.refresh = refresh_token
        self.refresh_url = refresh_url
        self.refresh_window = refresh_window
        self.expires = self._parse_expires_from_token(token)

    def login(self, username : str, password : str, url : str):
        req = Request(
            f"{url}/api/login",
            dumps({ "username" : username, "password": password }).encode("utf-8"),
            { "Content-Type" : "application/json" },
            method="POST"
        )
        self.refresh_url = f"{url}/api/refresh"
        self._request_token_update(req)

    def call_credentials(self):
        if self.refresh != "" and self.refresh_url != "" and datetime.now() >= self.expires - self.refresh_window :
            self.refresh_token()
        return access_token_call_credentials(self.token)

    def refresh_token(self):
        req = Request(
            self.refresh_url,
            dumps({ "refresh":self.refresh }).encode("utf-8"),
            { "Content-Type" : "application/json", "Authorization" : f"Bearer {self.token}" },
            method="POST"
        )
        self._request_token_update(req)

    def _parse_expires_from_token(self, token: str) -> datetime:
        claim_str = token.split(".")[1]
        claims = loads(urlsafe_b64decode(claim_str + ( '=' * ( 4 - len(claim_str) % 4 ) ) ))
        if "exp" in claims:
            return datetime.fromtimestamp(claims["exp"])
        return datetime.now() + timedelta(minutes=5)

    def _request_token_update(self, req : Request):
        with urlopen(req) as response:
            jsonRes = loads(response.read())
            self.token = jsonRes["token"]
            self.refresh = jsonRes["refresh"]
            self.expires = self._parse_expires_from_token(self.token)

