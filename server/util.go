package main

import (
	pb "blog-app/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
