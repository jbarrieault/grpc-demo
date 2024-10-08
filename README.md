# grpc-demo

Trying out gRPC, that's it.

## Generating gRPC services

The following generates the Echo service's gRPC service and protobuf code:

```shell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative services/echo/echo.proto
```

Start multiple instance of the echo service:
```shell
cd go-server
go run . -p 3000 &
go run . -p 3001 &
go run . -p 3002 &
```

The `go-server` also takes flags for simulating errors and latency. The following will start an EchoService server for which 5% of calls will produce an error, and 20% of (non-erroring) calls will be slow:
```
go run . -p 3000 -e 5 -s 20
```

Then start an echo service client, passing the addresses of the servers you started:
```shell
cd go-client && go run . -addr localhost:3000,localhost:3001,localhost:3002
```

Enter some messages and you will see responses coming from each server in round-robin fashion, skipping servers that it cannot establish connection to (due to invalid address, network partition, etc.)


## TODO

- [ ] auth?
- [ ] streaming server/client demo
- [ ] a ruby client
- [ ] demo proto schema evolution

