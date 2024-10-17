package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	pb "github.com/jbarrieault/grpc-demo/services/echo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:3100", "EchoService remote address as host:port")
)

func init() {
	flag.Parse()
	setupSignalHandler()
}

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(*addr, opts...)
	if err != nil {
		log.Fatalf("Failed to connect to Echo Service: %s", err)
	}
	defer conn.Close()

	// I now realize naming this the "streaming echo client" is misleading.
	// The client does not stream to the server. It's a client that makes a grpc
	// call for which the server streams a response.
	fmt.Println("Welcome to the Streaming Echo Client. Please enter a message...")

	client := pb.NewEchoClient(conn)
	message := pb.EchoNMessage{}
	message.N = 5 // hard-code for now
	reader := bufio.NewReader(os.Stdin)

	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")

		go echoN(input, &message, client)
	}
}

func echoN(input string, message *pb.EchoNMessage, client pb.EchoClient) {
	message.Value = input

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := client.EchoN(ctx, message)
	if err != nil {
		log.Println(err)
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			break
		}

		fmt.Printf("%s response %d/%d: %s\n", message.Value, msg.N+1, message.N, msg.Value)
	}
}

func setupSignalHandler() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nSee ya.")
		os.Exit(0)
	}()
}
