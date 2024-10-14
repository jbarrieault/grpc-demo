package main

import (
	"flag"

	pb "github.com/jbarrieault/grpc-demo/services/echo"
	"google.golang.org/grpc"
)

var (
	port = flag.String("p", "3100", "port to listen on")
)

type echoServer struct {
	pb.UnimplementedEchoServer
}

// EchoN streams the message.Value, message.N times,
// N is incremented on each successive streamed *EchoMessage.
func (c *echoServer) EchoN(message *pb.EchoNMessage, stream grpc.ServerStreamingServer[pb.EchoNMessage]) error {
	var err error
	for i := 0; i < int(message.N); i++ {
		m := &pb.EchoNMessage{N: int32(i), Value: message.Value}
		err = stream.Send(m)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	flag.Parse()
}
