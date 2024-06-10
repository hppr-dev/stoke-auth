from django.http import HttpResponse
from django.views.decorators.csrf import csrf_exempt
from stoke.django_annotations import require_token

def index(request):
    if request.user.is_anonymous:
        return HttpResponse(f"Hello unknown user!".encode())
    return HttpResponse(f"Hello {request.user}".encode())

@csrf_exempt
@require_token(claims={"inv" : "acc"})
def test(request):
    return HttpResponse(b'{ "hello": "world", "foo" : "bar" }', status=200, content_type="application/json")

