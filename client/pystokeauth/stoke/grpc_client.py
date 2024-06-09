from datetime import datetime, timedelta
from typing import Any, Dict
from grpc import access_token_call_credentials

from stoke.user_login import request_token, refresh_token, expires_from_token

class GRPCClientCredentials:
    def __init__(self, token : str, refresh_token : str = "", refresh_url : str = "", refresh_window : timedelta = timedelta(seconds=5)):
        self.token = token
        self.refresh = refresh_token
        self.refresh_url = refresh_url
        self.refresh_window = refresh_window
        self.expires = expires_from_token(token)

    def login(self, username : str, password : str, url : str):
        self.refresh_url = f"{url}/api/refresh"
        self._request_token_update(request_token(f"{url}/api/login", username, password))

    def call_credentials(self):
        if self.refresh != "" and self.refresh_url != "" and datetime.now() >= self.expires - self.refresh_window :
            self._request_token_update(refresh_token(self.refresh_url, self.token, self.refresh))
        return access_token_call_credentials(self.token)


    def _request_token_update(self, jsonRes : Dict[Any, Any]):
        self.token = jsonRes["token"]
        self.refresh = jsonRes["refresh"]
        self.expires = expires_from_token(self.token)

