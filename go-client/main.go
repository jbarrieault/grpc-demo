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
	addr = flag.String("addr", "localhost:3000", "The remote server address, as host:post")
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(*addr, opts...)
	if err != nil {
		log.Fatalf("Failed to connect to Echo Service: %s", err)
	}
	defer conn.Close()

	client := pb.NewEchoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nGoobye.")
		os.Exit(0)
	}()

	message := pb.EchoMessage{}
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to the Echo Client. Please enter a message...")

	for {
		input, _ := reader.ReadString('\n')
		message.Value = input

		resp, err := client.Echo(ctx, &message)
		if err != nil {
			log.Fatalf("Echo service error: %v", err)
		}

		log.Printf("Response: %v", resp.Value)
	}
}
