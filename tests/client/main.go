package main

import (
	"context"
	"log"

	"openmyth/go-kit.git/config"
	"openmyth/go-kit.git/pkg/grpc_client"
	"openmyth/go-kit.git/tests/helloworld"
)

func main() {
	client := grpc_client.NewGrpcClient(&config.Endpoint{
		Host: "0.0.0.0",
		Port: "8081",
	})
	ctx := context.Background()

	client.Connect(ctx)
	grpcClient := helloworld.NewGreeterClient(client.ClientConn)
	for range 5 {
		log.Println("1111")
		resp, err := grpcClient.SayHello(ctx, &helloworld.HelloRequest{
			Name: "hello",
		})
		if err != nil {
			panic(err)
		}
		log.Println(resp.GetMessage())
	}
	defer client.Close(ctx)
}
