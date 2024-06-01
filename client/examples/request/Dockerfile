FROM python:3.10-alpine3.20
RUN mkdir -p /app

WORKDIR /app

COPY . ./
RUN pip install -e ./pystokeauth

WORKDIR /app/examples/request
RUN pip install -r requirements.txt

EXPOSE 8080

ENTRYPOINT ["flask", "--app", "/app/examples/request/app.py", "run", "--port", "8080", "--host", "0.0.0.0"]

