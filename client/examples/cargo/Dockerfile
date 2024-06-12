FROM python:3.10-alpine3.20
RUN mkdir -p /app

WORKDIR /app

COPY . ./
RUN pip install -e ./pystokeauth
ENV CARGO_KEY=/app/examples/certs/cargo.key
ENV CARGO_CERT=/app/examples/certs/cargo.crt
ENV CARGO_CA=/app/examples/certs/ca.crt

WORKDIR /app/examples/cargo
RUN pip install -r requirements.txt

EXPOSE 6060

ENTRYPOINT ["python", "grpc_server.py"]

