package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	pb "github.com/jbarrieault/grpc-demo/services/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	port     = flag.String("p", "3000", "port to listen on")
	errRate  = flag.Int("e", 0, "percentage of calls that will error")
	slowRate = flag.Int("s", 0, "percentage of (non error) calls that be slow")

	kaep = keepalive.EnforcementPolicy{
		MinTime:             10 * time.Second, // terminate client connections if they ping more often than 10s
		PermitWithoutStream: true,             // Allow pings even when there are no active streams
	}
	kasp = keepalive.ServerParameters{
		MaxConnectionIdle: 10 * time.Second, // how long an idle connection can exist before terminating
	}
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

	s := grpc.NewServer(buildTlsConfig(), buildUserAuthenticationConfig(), grpc.KeepaliveEnforcementPolicy(kaep), grpc.KeepaliveParams(kasp))
	pb.RegisterEchoServer(s, &echoServer{})

	log.Printf("Echo Service listening on port %v", *port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func buildTlsConfig() grpc.ServerOption {
	crtPath, err := filepath.Abs("../tls/grpc-server.crt")
	if err != nil {
		log.Fatalf("failed to build crt file path: %v", err)
	}

	keyPath, err := filepath.Abs("../tls/grpc-server.key")
	if err != nil {
		log.Fatalf("failed to build key path: %v", err)
	}

	serverCert, err := tls.LoadX509KeyPair(crtPath, keyPath)
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}

	certPool := x509.NewCertPool()

	caPath, err := filepath.Abs("../tls/grpc-ca.crt")
	if err != nil {
		log.Fatalf("failed to build CA path: %v", err)
	}

	caCert, err := os.ReadFile(caPath)
	if err != nil {
		log.Fatalf("failed to read CA crt file: %v", err)
	}

	ok := certPool.AppendCertsFromPEM(caCert)
	if !ok {
		log.Fatalf("failed to add CA cert to the pool. Is %s a PEM file?", caPath)
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})

	return grpc.Creds(creds)
}

func buildUserAuthenticationConfig() grpc.ServerOption {
	return grpc.UnaryInterceptor(userAuthInterceptor)
}

func userAuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("userAuthInterceptor: no metadata provided")
		return nil, status.Errorf(codes.InvalidArgument, "missing authentication metadata")
	}

	vals := md.Get("authorization")
	if len(vals) == 0 {
		log.Println("userAuthInterceptor: no authorization metadadta provided")
		return nil, status.Errorf(codes.Unauthenticated, "missing authorization metadata")
	}

	jwt := strings.TrimPrefix(vals[0], "Bearer ")

	ok = validateJwt(jwt)
	if !ok {
		log.Printf("userAuthInterceptor: invalid token provided: '%v'", jwt)
		return nil, status.Errorf(codes.Unauthenticated, "invalid token provided: '%v'", jwt)
	}

	ret, err := handler(ctx, req)
	if err != nil {
		log.Printf("RPC failed: %v", err)
	}
	return ret, nil
}

func validateJwt(jwt string) bool {
	// using a hard-coded fake token
	return jwt == "MY.FAKE.JWT"
}
