package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	hw "seankhliao.com/grpctest/helloworld"
)

func main() {
	host := "grpc-server.seankhliao.com:8080"
	if len(os.Args) >= 2 {
		host = os.Args[1]
	}
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := hw.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &hw.HelloRequest{
		Name: "world",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply.GetMessage())
}
