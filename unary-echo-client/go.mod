module github.com/jbarrieault/grpc-demo/unary-echo-client

go 1.22.6

require (
	github.com/jbarrieault/grpc-demo/memory-registry v0.0.0-00010101000000-000000000000 // indirect
	github.com/jbarrieault/grpc-demo/services/echo v0.0.0-20241005113538-1fb53f7dfb29 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
	google.golang.org/grpc v1.67.1 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

replace github.com/jbarrieault/grpc-demo/services/echo => ../services/echo

replace github.com/jbarrieault/grpc-demo/memory-registry => ../memory-registry
