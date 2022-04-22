package main

import (
	pb "blog-app/proto"
	"bufio"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"log"
	"os"
	"strings"
)

const (
	create     = 0
	update     = 1
	get        = 2
	deleteById = 3
	getAll     = 4
	stop       = 5
)

func checkStdInErr(err error) {
	if err != nil {
		log.Fatalf("Error while reading input %v\n", err)
	}
}

func checkServerError(err error) {
	if err != nil {
		log.Fatalf("Server returned %v\n", err)
	}
}

func promptAction() int {
	choices := []string{"Create (c)", "Update (u)", "Get (g)", "Delete (d)", "GetAll (a)", "Stop (s)"}
	r := bufio.NewReader(os.Stdin)
	for {
		log.Printf("Choose which API to send request to, avilable are (%s)", choices)
		s, err := r.ReadString('\n')
		s = strings.ReplaceAll(s, "\n", "")
		checkStdInErr(err)
		switch s {
		case "c":
			return create
		case "u":
			return update
		case "g":
			return get
		case "d":
			return deleteById
		case "a":
			return getAll
		case "s":
			return stop
		default:
			return -1
		}
	}
}

func promptString(name string) string {
	log.Printf("-----> Enter %s:", name)
	r := bufio.NewReader(os.Stdin)
	s, err := r.ReadString('\n')
	s = strings.ReplaceAll(s, "\n", "")
	checkStdInErr(err)
	return s
}

func createBlog(client pb.BlogServiceClient) {
	a_ui := promptString("author id")
	title := promptString("title")
	content := promptString("content")
	m := pb.CreateBlogMessage{
		AuthorId: a_ui,
		Title:    title,
		Content:  content,
	}
	id, err := client.CreateBlog(context.Background(), &m)
	checkServerError(err)
	log.Printf("Created blog with id: %v\n", id)
}

func getBlogById(client pb.BlogServiceClient) {
	id := promptString("id")
	m := pb.BlogIdMessage{Id: id}
	blog, err := client.ReadBlog(context.Background(), &m)
	checkServerError(err)
	log.Printf("Here is your blog %v\n", blog)
}

func deleteBlogById(client pb.BlogServiceClient) {
	id := promptString("id")
	m := pb.BlogIdMessage{Id: id}
	_, err := client.DeleteBlog(context.Background(), &m)
	checkServerError(err)
	log.Println("Successfully delete blog")
}

func updateBlog(client pb.BlogServiceClient) {
	id := promptString("id")
	a_ui := promptString("author id")
	title := promptString("title")
	content := promptString("content")
	m := pb.BlogMessage{Id: id, AuthorId: a_ui, Title: title, Content: content}
	_, err := client.UpdateBlog(context.Background(), &m)
	checkServerError(err)
	log.Println("Successfully updated blog")
}

func listAllBlogs(client pb.BlogServiceClient) {
	stream, err := client.ListAllBlogs(context.Background(), &emptypb.Empty{})
	checkServerError(err)
	log.Printf("Here are all the blogs ----->")
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		checkServerError(err)
		log.Println(in)
	}
}

const ServerAddress = "localhost:30005"

func main() {
	connection, err := grpc.Dial(ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to %v\n", err)

	}
	client := pb.NewBlogServiceClient(connection)
	defer connection.Close()

	for {
		action := promptAction()
		switch action {
		case create:
			createBlog(client)
		case get:
			getBlogById(client)
		case deleteById:
			deleteBlogById(client)
		case update:
			updateBlog(client)
		case getAll:
			listAllBlogs(client)
		case stop:
			os.Exit(0)
		default:
			log.Printf("Unknow action")
		}
	}
}
