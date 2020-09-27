package main

import (
	"context"
	"log"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	pb "github.com/linuxxiaoyu/go-grpc-example/proto"
	"google.golang.org/grpc"
)

// type SearchService struct{}

// func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
// 	return &pb.SearchResponse{Response: r.GetRequest() + " Server"}, nil
// }

type MathService struct{}

func (s *MathService) Div(ctx context.Context, r *pb.MathRequest) (*pb.MathResponse, error) {
	return &pb.MathResponse{Ret: r.A / r.B}, nil
}

const PORT = "9001"

func main() {
	c, err := GetTLSCredentialsByCA()
	if err != nil {
		log.Fatalf("GetTLSCredentialsByCA err: %v", err)
	}

	opts := []grpc.ServerOption{
		grpc.Creds(c),
		grpc_middleware.WithUnaryServerChain(
			LoggingInterceptor,
			RecoveryInterceptor,
		),
	}

	server := grpc.NewServer(opts...)
	pb.RegisterMathServiceServer(server, &MathService{})

	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	server.Serve(lis)
}
