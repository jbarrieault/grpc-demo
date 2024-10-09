package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/jbarrieault/grpc-demo/services/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:3000", "Comma separated list of remote server(s) address, as host:post")
)

func init() {
	flag.Parse()
	registerStaticResolver()
	setupSignalHandler()
}

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`))

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
