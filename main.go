package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"regexp"

	pb "github.com/jschaf/bazel-go-gce"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func (s *server) Render(ctx context.Context, req *pb.RenderRequest) (*pb.RenderResponse, error) {
	italic := regexp.MustCompile(`\*(.*?)\*`)
	r := italic.ReplaceAllString(req.Text, "<em>$1</em>");
	return &pb.RenderResponse{Text: "<div>" + r + "</div>"}, nil
}

func main() {
	port := "8282"
	lis, err := net.Listen("tcp", ":" + port)
	defer lis.Close()
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMarkdownServer(s, &server{})
	reflection.Register(s)

	fmt.Println("===========================")
	fmt.Println("Server started on port " + port)
	fmt.Println("===========================")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
