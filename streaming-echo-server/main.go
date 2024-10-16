package main

import (
	"flag"
	"log"
	"net"
	"time"

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

		if i < int(message.N)-1 {
			time.Sleep(time.Duration(1) * time.Second)
		}
	}

	return nil
}

func main() {
	flag.Parse()

	listner, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("net.Listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterEchoServer(s, &echoServer{})

	log.Printf("Streaming Echo Server listening on port %v", *port)
	err = s.Serve(listner)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
