# grpc-demo

Trying out gRPC, that's it.

## Generating gRPC services

The following generates the Echo service's gRPC service and protobuf code:

`protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative services/echo/echo.proto`

```shell
# start a echo service server(s)
cd go-server && go run . -p 3000 &&
cd go-server && go run . -p 3001 &&
cd go-server && go run . -p 3002 &&

# start a client, passing the addresses of the servers you started
cd go-client && go run . -addr localhost:3000,localhost:3001,localhost:3002

# enter some messages and you will see responses coming from each server in round-robin fashion
```
