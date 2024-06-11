from typing import Dict, Optional
from django.contrib.auth.models import User
from django.contrib.auth.backends import ModelBackend
from django.core.exceptions import PermissionDenied
from django.http import HttpRequest
from django.conf import settings

from stoke.client import StokeClient
from stoke.test_client import TestStokeClient
from stoke.user_login import request_token

class StokeAuthBackend(ModelBackend):
    """
        Auth backend to use a stoke server to authenticate
        Requires the following settings:
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
    """

    def authenticate(self, request : HttpRequest, username : Optional[str] = None, password : Optional[str] = None, claims : Optional[Dict[str, str]] = None):
        conf = settings.STOKE_AUTH_CONFIG
        client : StokeClient | TestStokeClient = conf["CLIENT"]

        if request.user is not None and not request.user.is_anonymous:
            return request.user

        if type(client) is TestStokeClient:
            token = client.default_token
        elif "HTTP_AUTHORIZATION" in request.META and request.META["HTTP_AUTHORIZATION"].startswith("Bearer "):
            token = request.META["HTTP_AUTHORIZATION"].removeprefix("Bearer ")
        elif username is not None and password is not None:
            # If username and password are supplied, we need to check the token.
            token = self._request_token(client, username, password)
            if claims is None:
                # Set claims to empty by default to allow django to do the permissions when user/password
                claims = {}
            if token == "" :
                return None
        else:
            return None

        if claims is None:
            if "ROUTE_CLAIMS" in conf and request.path in conf["ROUTE_CLAIMS"]:
                claims = conf["ROUTE_CLAIMS"]
            else:
                claims = conf["BASE_CLAIMS"]


        if token is None:
            raise PermissionDenied()

        jwtDict = client.parse_token(token)
        if jwtDict is None:
            raise PermissionDenied()

        for k, v in claims.items():
            if k not in jwtDict or jwtDict[k] != v:
                raise PermissionDenied()

        request.META["STOKE_AUTH_CLAIMS"] = jwtDict
        request.META["STOKE_AUTH_TOKEN"] = token

        try:
            user = User.objects.get(username=jwtDict[conf["USERNAME"]])
        except User.DoesNotExist:
            fname = "" if "FIRST_NAME" not in conf else jwtDict[conf["FIRST_NAME"]]
            lname = "" if "LAST_NAME" not in conf else jwtDict[conf["LAST_NAME"]]
            email = "" if "EMAIL" not in conf else jwtDict[conf["EMAIL"]]
            user = User(
                username = jwtDict[conf["USERNAME"]],
                first_name = fname,
                last_name = lname,
                email = email,
            )
            user.is_superuser = "SUPERUSER" in conf and  jwtDict.get(conf["SUPERUSER"][0]) == conf["SUPERUSER"][1] 
            user.is_staff = "STAFF" in conf and jwtDict.get(conf["STAFF"][0]) == conf["STAFF"][1]
            user.save()
        return user

    def _request_token(self, client : StokeClient | TestStokeClient, username : str, password : str) -> str:
        resp = request_token(f"{client.url}/api/login", username, password)
        return resp["token"]

