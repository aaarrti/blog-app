package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func _connect() *mongo.Collection {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("blog-db").Collection("blogs")
	log.Printf("Connected successfully to MongoDB")
	return collection
}

func save(ctx context.Context, item *CreateBlogItem) (*string, error) {
	var collection = _connect()
	res, err := collection.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}
	iod, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Error(codes.Internal, "Failed to convert to OId")
	}
	str := iod.Hex()
	return &str, nil
}

func update(ctx context.Context, item *BlogItem) error {
	var collection = _connect()
	res, err := collection.UpdateOne(ctx,
		bson.M{"_id": item.ID},
		bson.M{"$set": item},
	)
	if err != nil {
		log.Printf("Failed to update document %v\n", item.ID)
		return status.Errorf(codes.Internal, "Could not update")
	}
	if res.MatchedCount == 0 {
		log.Printf("Matched count == 0")
		return status.Errorf(codes.NotFound, "Blog not found")
	}
	return nil
}
