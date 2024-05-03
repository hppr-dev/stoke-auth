from typing import Dict
from functools import wraps
from flask import request
from werkzeug.exceptions import Unauthorized

from stoke.client import StokeClient

def require_claims(client : StokeClient, claims : Dict[str, str], jwt_kwarg : str | None = None):
    def inner(func):
        @wraps(func)
        def wrap_require(*args, **kwargs):
            if request.authorization is None:
                raise Unauthorized(description="Missing Authorization Header")

            if request.authorization.type != "bearer" or request.authorization.token is None:
                raise Unauthorized(description="Missing Authorization Token")

            jwtDict = client.parse_token(request.authorization.token)
            if jwtDict is None:
                raise Unauthorized(description="Invalid Token")

            for k, v in claims.items():
                if k not in jwtDict or jwtDict[k] != v:
                    raise Unauthorized(description="Missing required claims")

            if jwt_kwarg is not None:
                kwargs[jwt_kwarg] = jwtDict
            return func(*args, **kwargs)

        return wrap_require
    return inner

