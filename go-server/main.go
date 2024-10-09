package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	pb "github.com/jbarrieault/grpc-demo/services/echo"
	"google.golang.org/grpc"
)

var (
	port = flag.String("p", "3000", "port to listen on")
)

type echoServer struct {
	pb.UnimplementedEchoServer
}

func (*echoServer) Echo(_ context.Context, message *pb.EchoMessage) (*pb.EchoMessage, error) {
	log.Printf("Received Echo: %v", message.GetValue())
	simulatedLatency := rand.Intn(6)
	if simulatedLatency >= 5 {
		fmt.Println("simulating a network hiccup causing increased latency...")
		time.Sleep(time.Duration(simulatedLatency) * time.Second)
	}

	respMsg := fmt.Sprintf("Server at port %v: %v", *port, message.GetValue())
	return &pb.EchoMessage{Value: respMsg}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("net.Listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterEchoServer(s, &echoServer{})

	log.Printf("Echo Service listening on port %v", *port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
