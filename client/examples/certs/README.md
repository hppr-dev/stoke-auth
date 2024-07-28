# Example certificates directory

This directory is for running the example environment.
Procedures to generate/clean this data are implemented in the `t cert` task

## Certificates

The following certificates must be generated to use/test the example environment:

  * ca.key/ca.crt -- The ca certificate
  * cargo.key/cargo.crt -- the cargo hold grpc certificate for use with mTLS

## Generating certificates with openssl

Each certificate/key pair has an associated config file in the config directory.

To generate a cert/key pair run the following command (remove -CA/-CAkey for the ca):
```
    openssl req -x509 -newkey rsa:2048 -config config/NAME.conf -out NAME.crt -CA ca.crt -CAkey ca.key 
```

The generated certificates are environment specific and should be able to be removed/regenerated as needed.

Use `openssl verify -CAfile ca.crt *.crt` to verify generated certificates
