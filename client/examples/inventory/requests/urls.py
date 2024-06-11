from django.urls import path

from . import views

urlpatterns = [
    path("", views.index, name="index"),
    path("test/", views.test, name="test"),
    path("cargo_contents/", views.cargo_contents, name="cargo_contents"),
]
