package main

import (
	"context"
	"log"
	"net"

	pb "github.com/jbarrieault/grpc-demo/services/echo"
	"google.golang.org/grpc"
)

type echoServer struct {
	pb.UnimplementedEchoServer
}

func (*echoServer) Echo(_ context.Context, message *pb.EchoMessage) (*pb.EchoMessage, error) {
	log.Printf("Received Echo: %v", message.GetValue())
	return &pb.EchoMessage{Value: message.GetValue()}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("net.Listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterEchoServer(s, &echoServer{})

	log.Println("Echo Service listening on port 3000")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
