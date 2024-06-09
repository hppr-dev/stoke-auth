from django.http import HttpResponse
from stoke.django_annotations import require_token

@require_token(claims={"inv" : "acc"}, start_session=True)
def index(request):
    return HttpResponse(f"Hello {request.user}".encode())

@require_token(claims={"inv" : "acc"}, start_session=True)
def test(request):
    return HttpResponse(f"Hello {request.user}".encode())
