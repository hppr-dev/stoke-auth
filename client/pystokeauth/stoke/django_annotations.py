from typing import Dict, Optional
from functools import wraps
from django.http import HttpRequest, HttpResponse
from django.contrib.auth import authenticate, login

def require_token(claims : Optional[Dict[str, str]] = None, start_session : bool = False):
    """
        Requires a stoke token with certain claims present to access
        If claims are ommited, the ROUTE_CLAIMS or BASE_CLAIMS settings will be used.
        Set start_session = true to start a django session to have django to handle auth after.
        Claim information is not stored when transfering to a django session.
    """
    def inner(func):
        @wraps(func)
        def wrap_require(*args, **kwargs):
            req : HttpRequest = args[0]

            user = authenticate(req, claims=claims)
            if user is None:
                return HttpResponse('Unauthorized'.encode(), status=401)

            if start_session:
                login(req, user)

            return func(*args, **kwargs)

        return wrap_require
    return inner
