FROM python:3.10-alpine3.20
RUN mkdir -p /app

WORKDIR /app

COPY . ./
RUN pip install -e ./pystokeauth

WORKDIR /app/examples/cargo
RUN pip install -r requirements.txt

EXPOSE 6060

ENTRYPOINT ["python", "grpc_server.py"]

