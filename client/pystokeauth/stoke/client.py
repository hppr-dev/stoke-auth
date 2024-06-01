from datetime import datetime, timezone
from typing import Dict, Optional, Any
from ssl import SSLContext

from base64 import urlsafe_b64decode
from jwt import PyJWK, PyJWKClient, PyJWKSet, decode
from jwt.types import JWKDict

from cryptography.hazmat.primitives.asymmetric.rsa import RSAPublicNumbers

class StokeClient(PyJWKClient):
    def __init__(self, url: str, headers : Optional[Dict[str, Any]] = None, ssl_context : Optional[SSLContext] = None):
        super().__init__(
            uri = url,
            cache_keys = True,
            max_cached_keys = 2,
            cache_jwk_set = True,
            headers = headers,
            timeout = 30,
            ssl_context = ssl_context
        )
        self.next_update : datetime = datetime.now(timezone.utc)

    def get_jwk_set(self, refresh: bool = False) -> PyJWKSet:
        return super().get_jwk_set(refresh or datetime.now(timezone.utc) >= self.next_update)

    def fetch_data(self) -> Any:
        data = super().fetch_data()
        self.next_update = datetime.strptime(data["exp"], "%Y-%m-%dT%H:%M:%S%z")
        return data

    def parse_token(self, token : str) -> Optional[Dict[str, Any]]:
        keys = self.get_signing_keys()
        for key in keys:
            try:
                decoded = decode(token, key.key, algorithms=[self.__get_algo_str(key)])
                return decoded
            except Exception as e:
                print(e)
        return None

    def __get_algo_str(self, data: PyJWK) -> str :
        jwk = data._jwk_data
        if jwk["kty"] == "EC":
            if jwk["crv"] == "P-256":
                return "ES256"
            if jwk["crv"] == "P-384":
                return "ES384"
            if jwk["crv"] == "P-521":
                return "ES512"
        if jwk["kty"] == "RSA":
            return self.__get_rsa_algo(jwk)
        if jwk["kty"] == "OKP":
            return "EdDSA"
        return "Unknown"

    def __get_rsa_algo(self, jwk : JWKDict) -> str :
        n = int.from_bytes(urlsafe_b64decode(jwk["n"]), "big")
        e = int.from_bytes(urlsafe_b64decode(jwk["e"]), "big")
        size = RSAPublicNumbers(n, e).public_key().key_size
        if size == 256:
            return "PS256"
        if size == 384:
            return "PS384"
        if size == 512:
            return "PS521"
        return "Unknown"

        
