from django.urls import re_path

from requests import sockets

websocket_urlpatterns = [
    re_path(r"load_cargo/$", sockets.LoadItemConsumer.as_asgi()),
]
