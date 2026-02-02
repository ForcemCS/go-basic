package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "my-grpc-demo/pb"
)

func main() {
	// 使用 insecure 是因为本地测试不想配 SSL 证书
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "小白"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Printf("Server Response: %s\n", r.GetReply())
}
