# grpc-demo

Trying out gRPC, that's it.

> [!TIP]
> These experimental programs resemble work in the [official examples](https://github.com/grpc/grpc-go/tree/master/examples).
> If you want to learn gRPC, look there instead.


## Generating gRPC services

The following re-generates the Echo service's gRPC service and protobuf code:

```shell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative services/echo/echo.proto
```

### Unary Echo

The unary echo client/server demo a simple call-response model, configured with some additional gRPC feaures:
- Load-balancing (round robin)
- A logging & response decorating interceptor
- A user authentication interceptor
- A Service address resolver
- A back-off retry policy
- mTLS
- Keepalive pings

Be sure to follow the [instructions](./tls/README.md) to set up the required keys & SSL certificates.

Start multiple instance of the echo service:
```shell
cd unary-echo-server
go run . -p 3000 &
go run . -p 3001 &
go run . -p 3002 &
```

The `unary-echo-server` also takes flags for simulating errors and latency. The following will start an EchoService server for which 5% of calls will produce an error, and 20% of (non-erroring) calls will be slow:
```
go run . -p 3000 -e 5 -s 20
```

Then start an echo service client, passing the addresses of the servers you started:
```shell
cd unary-echo-client && go run . -addr localhost:3000,localhost:3001,localhost:3002
```

Enter some messages and you will see responses coming from each server in round-robin fashion, skipping servers that it cannot establish connection to (due to invalid address, network partition, etc.)

### Server Streaming Echo

The streaming echo example showcases gRPC server-streaming.
In this case, the server responds to single client request with a stream of messages.

```shell
cd streaming-echo-server
go run .
```
```shell
cd streaming-echo-client
go run .
```

When the client sends a message, you will see a stream of responses.
Multiple messages can be sent in short succession, which results in the interleaving of response stream messages.

<img width="616" alt="image" src="https://github.com/user-attachments/assets/9c405969-3fdf-4dd7-865d-5926978bd91e">

## TODO

- [X] build an (in-memory) service registry module
- [X] create a non-static resolver using a service registry
- [X] explore TLS/mTLS
- [X] token based user auth interceptor
- [X] streaming server demo
- [X] keepalive
- [/] health checking
- [ ] streaming client demo
- [ ] streaming bi-di demo
- [ ] observability/metrics
- [ ] a ruby client
- [ ] demo proto schema evolution
