# Encryption Setup

Some examples enable TLS and mTLS encryption, which requires generating
keys and certificates.

First, generate a self-signed CA cert:
```shell
# generate a private key for the CA
openssl genrsa -out grpc-ca.key 2048

# generate a self-signed CA certificate
openssl req -x509 -new -nodes -key grpc-ca.key -sha256 -days 365 -out grpc-ca.crt -subj "/CN=gRPC Demo CA"
```

Then, generate a key & cert signing request for the gRPC server and client:
```shell
openssl genrsa -out grpc-server.key 2048
openssl req -new -key grpc-server.key -out grpc-server.csr -config grpc-server.cnf

openssl genrsa -out grpc-client.key 2048
openssl req -new -key grpc-client.key -out grpc-client.csr -config grpc-client.cnf
```

Finally, generate the server's cert by having the CA sign the CSR:
```shell
openssl x509 -req -in grpc-server.csr -CA grpc-ca.crt -CAkey grpc-ca.key -CAcreateserial -out grpc-server.crt -days 365 -sha256 -extfile grpc-server.cnf -extensions req_ext
openssl x509 -req -in grpc-client.csr -CA grpc-ca.crt -CAkey grpc-ca.key -CAcreateserial -out grpc-client.crt -days 365 -sha256 -extfile grpc-client.cnf -extensions req_ext
```
