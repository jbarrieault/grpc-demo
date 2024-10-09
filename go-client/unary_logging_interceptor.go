package main

import (
	"context"
	"fmt"
	"log"
	"os"

	pb "github.com/jbarrieault/grpc-demo/services/echo"
	"google.golang.org/grpc"
)

func unaryLoggingInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	reqMessage, ok := req.(*pb.EchoMessage)

	if ok {
		log.Printf("unaryLoggingInterceptor pre-processing: %v", reqMessage.Value)
	} else {
		fmt.Fprint(os.Stderr, "req was not *pb.EchoMessage")
	}

	err := invoker(ctx, method, req, reply, cc, opts...)

	// TODO: log which host served the request

	if err != nil {
		return err
	}

	resMessage, ok := reply.(*pb.EchoMessage)
	if !ok {
		fmt.Fprint(os.Stderr, "res was not *pb.EchoMessage")
		return nil
	}

	log.Printf("unaryLoggingInterceptor post-processing: decorating res \"%v\")", resMessage.Value)

	resMessage.Value = fmt.Sprintf("<decorated>%v</decorated>", resMessage.Value)

	return err
}
