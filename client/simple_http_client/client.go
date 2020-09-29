package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	pb "github.com/linuxxiaoyu/go-grpc-example/proto"
	"google.golang.org/grpc"

	"github.com/linuxxiaoyu/go-grpc-example/pkg/gtls"
)

const PORT = "9003"

func main() {
	tlsClient := gtls.Client{
		ServerName: "go-grpc-example",
		CertFile:   "../../conf/server/server.pem",
	}
	c, err := tlsClient.GetTLSCredentials()
	if err != nil {
		log.Fatalf("tlsClient.GetTLSCredentials err: %v", err)
	}

	conn, err := grpc.Dial(":"+PORT, grpc.WithTransportCredentials(c))
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := pb.NewSearchServiceClient(conn)
	resp, err := client.Search(ctx, &pb.SearchRequest{
		Request: "gRPC",
	})
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Fatalln("client.Search err: deadline")
			}
		}
		log.Fatalf("client.Search err: %v", err)
	}

	log.Printf("resp: %s", resp.GetResponse())
}
