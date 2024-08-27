package main

import (
	"context"
	v1 "example-grpc/gen/proto/hello"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	v1.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	return &v1.HelloReply{Mesasge: "Hello " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	v1.RegisterGreeterServer(s, &server{})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
