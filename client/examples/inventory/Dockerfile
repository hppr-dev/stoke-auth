FROM python:3.10-alpine3.20
RUN mkdir -p /app

WORKDIR /app

COPY . ./
RUN pip install -e ./pystokeauth
RUN pip install -e ./examples/cargo

ENV CA_FILE=/app/examples/certs/ca.crt
ENV KEY_FILE=/app/examples/certs/inventory.key
ENV CERT_FILE=/app/examples/certs/inventory.crt
ENV STOKE_ENV=docker

WORKDIR /app/examples/inventory
RUN pip install -r requirements.txt
RUN python manage.py migrate

EXPOSE 8080

# Don't do this
ENTRYPOINT ["python", "manage.py", "runserver", "0.0.0.0:8080"]

