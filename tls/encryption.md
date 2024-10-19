# Encryption Setup

Some examples enable TLS and mTLS encryption, which requires generating
certificates and keys.


First, generate a self-signed CA cert:
```
# generate a private key for the CA
openssl genpkey -algorithm RSA -out grpc-ca.key -pkeyopt rsa_keygen_bits:2048

# generate a self-signed CA certificate
openssl req -x509 -new -nodes -key grpc-ca.key -sha256 -days 365 -out grpc-ca.crt -subj "/CN=gRPC Demo CA"
```

Then, generate a key & cert signing request for the gRPC server:
```
openssl genpkey -algorithm RSA -out grpc-server.key -pkeyopt rsa_keygen_bits:2048
openssl req -new -key grpc-server.key -out grpc-server.csr -subj "/CN=gRPC Demo Server"
```

Finally, generate the server's cert by having the CA sign the CSR:
```
openssl x509 -req -in grpc-server.csr -CA grpc-ca.crt -CAkey grpc-ca.key -CAcreateserial -out grpc-server.crt -days 365 -sha256
```
