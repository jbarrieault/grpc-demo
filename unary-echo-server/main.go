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
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	port     = flag.String("p", "3000", "port to listen on")
	errRate  = flag.Int("e", 0, "percentage of calls that will error")
	slowRate = flag.Int("s", 0, "percentage of (non error) calls that be slow")

	kaep = keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second, // terminate client connections if they ping more often than 10s
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}

	// Note about keepalive: Ping activity only acts to keep _idle_ connections alive.
	// A client is considered idle based on RPC call activity, which does not include ping activity.
	kasp = keepalive.ServerParameters{
		MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
		MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
		MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
		Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
		Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
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
	healthCheckServer := health.NewServer()
	healthgrpc.RegisterHealthServer(s, healthCheckServer)
	pb.RegisterEchoServer(s, &echoServer{})

	// after 5 seconds become unhealthy for 10 seconds, then become healthy again
	go func() {
		time.Sleep(5 * time.Second)
		system := "" // I believe this indicates _all_ registered services.
		healthCheckServer.SetServingStatus(system, healthgrpc.HealthCheckResponse_NOT_SERVING)
		log.Printf("System health status changed to %v", healthgrpc.HealthCheckResponse_NOT_SERVING)

		time.Sleep(10 * time.Second)
		healthCheckServer.SetServingStatus(system, healthgrpc.HealthCheckResponse_SERVING)
		log.Printf("System health status changed to %v", healthgrpc.HealthCheckResponse_SERVING)
	}()

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
