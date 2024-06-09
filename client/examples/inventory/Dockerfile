FROM python:3.10-alpine3.20
RUN mkdir -p /app

WORKDIR /app

COPY . ./
RUN pip install -e ./pystokeauth

WORKDIR /app/examples/inventory
RUN pip install -r requirements.txt
RUN python manage.py migrate

EXPOSE 8080

# Don't do this
ENTRYPOINT ["python", "manage.py", "runserver", "0.0.0.0:8080"]

