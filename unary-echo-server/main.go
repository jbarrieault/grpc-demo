package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"path/filepath"
	"time"

	pb "github.com/jbarrieault/grpc-demo/services/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

var (
	port     = flag.String("p", "3000", "port to listen on")
	errRate  = flag.Int("e", 0, "percentage of calls that will error")
	slowRate = flag.Int("s", 0, "percentage of (non error) calls that be slow")
)

type echoServer struct {
	pb.UnimplementedEchoServer
}

func (*echoServer) Echo(_ context.Context, message *pb.EchoMessage) (*pb.EchoMessage, error) {
	log.Printf("Received Echo: %v", message.GetValue())

	if *errRate > 0 && *errRate >= rand.Intn(100) {
		fmt.Println("Simulating an UNAVAILABLE error")
		return nil, status.Errorf(codes.Unavailable, "Service unavailable (simulated)")
	}

	if *slowRate > 0 && *slowRate >= rand.Intn(100) {
		fmt.Println("Simulating slow request")
		time.Sleep(time.Duration(3) * time.Second)
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

	cert, err := filepath.Abs("../tls/grpc-server.crt")
	if err != nil {
		log.Fatalf("failed to build cert path: %v", err)
	}

	pkey, err := filepath.Abs("../tls/grpc-server.key")
	if err != nil {
		log.Fatalf("failed to build private key path: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile(cert, pkey)
	if err != nil {
		log.Fatalf("failed to create server credentials: %v", err)
	}

	credsServerOpt := grpc.Creds(creds)

	s := grpc.NewServer(credsServerOpt)
	pb.RegisterEchoServer(s, &echoServer{})

	log.Printf("Echo Service listening on port %v", *port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
