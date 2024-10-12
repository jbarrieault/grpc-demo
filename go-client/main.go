package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	// mr "github.com/jbarrieault/grpc-demo/go-client/pkg/memory_registry"

	pb "github.com/jbarrieault/grpc-demo/services/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr          = flag.String("addr", "localhost:3000", "Comma separated list of remote server(s) address, as host:post")
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
)

func init() {
	flag.Parse()
	registerStaticResolver()
	setupSignalHandler()
}

func main() {
	// r := mr.NewRegistery()
	// err := r.Register("echo", strings.Split(*addr, ",")...)
	// if err != nil {
	// 	log.Println(err)
	// }

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithUnaryInterceptor(unaryLoggingInterceptor))
	opts = append(opts, grpc.WithDefaultServiceConfig(serviceConfig))

	conn, err := grpc.NewClient("static:///i-believe-this-is-ignored-and-the-resolver-takes-over?", opts...)
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

func echo(input string, message *pb.EchoMessage, client pb.EchoClient) (string, error) {
	message.Value = input

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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
