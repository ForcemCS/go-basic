package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "my-grpc-demo/pb"
)

type server struct {
	pb.UnimplementedGreeterServer
}

// 2. 实现 SayHello 方法
// 注意：这个函数的签名是生成的代码里规定死的，我们必须照着写
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Reply: "Hello " + req.GetName()}, nil
}

func main() {
	// 监听 50051 端口
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 创建 gRPC 服务器
	s := grpc.NewServer()
	// 注册服务
	// RegisterGreeterServer 这个函数就是生成的代码里提供的！
	// 它把我们的 server 结构体注册进了 gRPC 内部
	pb.RegisterGreeterServer(s, &server{})

	fmt.Println("Server listening at :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
