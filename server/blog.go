package main

import (
	pb "blog-app/proto"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CreateBlogItem struct {
	AuthorId string `bson:"author_id"`
	Title    string `bson:"title"`
	Content  string `bson:"content"`
}

type BlogItem struct {
	ID       primitive.ObjectID `bson:"_id"`
	AuthorId string             `bson:"author_id"`
	Title    string             `bson:"title"`
	Content  string             `bson:"content"`
}

func toMessage(data *BlogItem) *pb.BlogMessage {
	return &pb.BlogMessage{
		Id:       data.ID.Hex(),
		AuthorId: data.AuthorId,
		Title:    data.Title,
		Content:  data.Content,
	}
}

func toCreateEntity(m *pb.CreateBlogMessage) *CreateBlogItem {
	return &CreateBlogItem{
		AuthorId: m.AuthorId,
		Title:    m.Title,
		Content:  m.Content,
	}
}

func toEntity(m *pb.BlogMessage) (*BlogItem, error) {
	id, err := primitive.ObjectIDFromHex(m.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Failed to parse Id")
	}
	return &BlogItem{
		ID:       id,
		AuthorId: m.AuthorId,
		Title:    m.Title,
		Content:  m.Content,
	}, nil
}

func (s *Server) CreateBlog(c context.Context, message *pb.CreateBlogMessage) (*pb.BlogIdMessage, error) {
	entity := toCreateEntity(message)
	res, err := save(c, entity)
	if err != nil {
		return nil, err
	}
	return &pb.BlogIdMessage{Id: *res}, nil
}

func (s *Server) UpdateBlog(c context.Context, message *pb.BlogMessage) (*emptypb.Empty, error) {
	entity, err := toEntity(message)
	if err != nil {
		return nil, err
	}
	err = update(c, entity)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) ReadBlog(c context.Context, in *pb.BlogIdMessage) (*pb.BlogMessage, error) {
	res, err := findById(c, in.Id)
	if err != nil {
		return nil, err
	}
	message := toMessage(res)
	return message, nil
}
