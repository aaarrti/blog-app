package main

import (
	pb "blog-app/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

type Server struct {
	pb.BlogServiceServer
}

func main() {
	port, ok := os.LookupEnv("SERVER_PORT")
	if !ok {
		log.Fatalf("SERVER_PORT env var not set")
	}
	addr := "0.0.0.0:" + port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on :%v\n", err)
	}
	log.Printf("Listeninng on %s\n", addr)

	server := grpc.NewServer()
	pb.RegisterBlogServiceServer(server, &Server{})
	reflection.Register(server)
	if err = server.Serve(listener); err != nil {
		log.Fatalf("Failed to server: %v\n", err)
	}

}
