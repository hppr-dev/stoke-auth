from typing import Any, Dict, Optional
from jwt import decode
from stoke.client import StokeClient


class TestStokeClient(StokeClient):
    def __init__(self, default_dict : Dict[str, Any] = {}, default_token : str = ""):
        self.default_dict = default_dict
        self.default_token = default_token
        self.url = None

    def parse_token(self, token: str) -> Optional[Dict[str, Any]]:
        if len(self.default_dict) > 0:
            return self.default_dict
        return decode(token, options={"verify_signature": False})

