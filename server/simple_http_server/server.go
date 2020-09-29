package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/linuxxiaoyu/go-grpc-example/proto"
	"google.golang.org/grpc"

	"github.com/linuxxiaoyu/go-grpc-example/pkg/gtls"
)

type SearchService struct{}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			return nil, status.Errorf(codes.Canceled, "SearchService.Search canceled")
		}

		time.Sleep(time.Second)
	}

	return &pb.SearchResponse{Response: r.GetRequest() + " HTTP Server"}, nil
}

const PORT = "9003"

func main() {
	certFile := "../../conf/server/server.pem"
	keyFile := "../../conf/server/server.key"
	tlsServer := gtls.Server{
		CertFile: certFile,
		KeyFile:  keyFile,
	}
	c, err := tlsServer.GetTLSCredentials()
	if err != nil {
		log.Fatalf("tlsServer.GetTLSCredentials err: %v", err)
	}

	mux := GetHTTPServeMux()
	server := grpc.NewServer(grpc.Creds(c))
	pb.RegisterSearchServiceServer(server, &SearchService{})

	http.ListenAndServeTLS(
		":"+PORT,
		certFile,
		keyFile,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
				server.ServeHTTP(w, r)
			} else {
				mux.ServeHTTP(w, r)
			}
			return
		}),
	)
}

func GetHTTPServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("linuxxiaoyu: go-grpc-example"))
	})

	return mux
}
