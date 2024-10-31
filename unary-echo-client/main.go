package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	mr "github.com/jbarrieault/grpc-demo/memory-registry"
	pb "github.com/jbarrieault/grpc-demo/services/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var (
	addr          = flag.String("addr", "localhost:3000", "Comma separated list of remote server(s), as host:port")
	serviceConfig = `{
			"methodConfig": [{
				"name": [{"service": "echo.Echo"}],
				"retryPolicy": {
					"MaxAttempts": 4,
					"InitialBackoff": "1.00s",
					"MaxBackoff": "20.0s",
					"BackoffMultiplier": 2,
					"RetryableStatusCodes": [ "UNAVAILABLE" ]
				}
			}],
			"loadBalancingConfig": [{"round_robin":{}}]
		}`
	serverName = "server.grpc-demo.example.com"
	mem_reg    *mr.MemoryRegistry
)

func init() {
	flag.Parse()
	setupSignalHandler()
	initMemoryRegistry()
	registerRegistryResolver()
	registerStaticResolver()
}

func main() {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(buildTlsConfig()))
	opts = append(opts, grpc.WithUnaryInterceptor(unaryLoggingInterceptor))
	opts = append(opts, grpc.WithDefaultServiceConfig(serviceConfig))

	// 'static://' scheme is handled by static_resolver, which uses a hard-coded list of addresses
	// conn, err := grpc.NewClient("static:///this-part-doesnt-matter-because-the-static-resolver-is-static", opts...)

	// the registry schema is handled by registry_resolver, which looks up address from a an in-memory service registry
	conn, err := grpc.NewClient("registry:///echo.Echo", opts...)
	if err != nil {
		log.Fatalf("Failed to connect to Echo Service: %s", err)
	}
	defer conn.Close()

	client := pb.NewEchoClient(conn)
	message := pb.EchoMessage{}
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to the Echo Client. Please enter a message...")

	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")
		output, err := echo(input, &message, client)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		fmt.Println(output)
	}
}

func buildTlsConfig() credentials.TransportCredentials {
	crtPath, err := filepath.Abs("../tls/grpc-client.crt")
	if err != nil {
		log.Fatalf("failed to build crt file path: %v", err)
	}

	keyPath, err := filepath.Abs("../tls/grpc-client.key")
	if err != nil {
		log.Fatalf("failed to build key path: %v", err)
	}

	cert, err := tls.LoadX509KeyPair(crtPath, keyPath)
	if err != nil {
		log.Fatalf("failed to load client cert: %v", err)
	}

	certPool := x509.NewCertPool()

	caCrtPath, err := filepath.Abs("../tls/grpc-ca.crt")
	if err != nil {
		log.Fatalf("failed to build CA path: %v", err)
	}

	caCrt, err := os.ReadFile(caCrtPath)
	if err != nil {
		log.Fatalf("failed to read CA crt file: %v", err)
	}

	ok := certPool.AppendCertsFromPEM(caCrt)
	if !ok {
		log.Fatalf("failed to add CA cert to the pool. Is %s a PEM file?", caCrtPath)
	}

	tlsConfig := &tls.Config{
		ServerName:   serverName,
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(tlsConfig)
}

func echo(input string, message *pb.EchoMessage, client pb.EchoClient) (string, error) {
	message.Value = input

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// using a hard-coded fake token
	jwt := metadata.Pairs("authorization", "Bearer MY.FAKE.JWT")
	ctx = metadata.NewOutgoingContext(ctx, jwt)

	resp, err := client.Echo(ctx, message)
	if err != nil {
		return "", err
	}

	return resp.Value, nil
}

func setupSignalHandler() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nGoodbye.")
		os.Exit(0)
	}()
}

// TODO: It would be fun to extract the registry
// to its own process, exposed over a socket.
// the unary-echo-server program could register itself on startup,
// which would much closer to a real-world setup.
func initMemoryRegistry() {
	mem_reg = mr.NewRegistery()
	err := mem_reg.Register("echo.Echo", strings.Split(*addr, ",")...)
	if err != nil {
		log.Println(err)
	}
}
