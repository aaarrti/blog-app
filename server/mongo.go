package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"time"
)

var client *mongo.Client = func() *mongo.Client {
	mongoUri, ok := os.LookupEnv("MONGO_URI")
	if !ok {
		log.Fatalf("MONGO_URI env var not set")
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected successfully to MongoDB")
	return client
}()

func testConnection() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Failed to ping MongoDB")
	}
}

func disconnect(client *mongo.Client) {
	if err := client.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}

func getCollection(client *mongo.Client) *mongo.Collection {
	err := client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	_collection := client.Database("blog-db").Collection("blogs")
	return _collection
}

func save(ctx context.Context, item *CreateBlogItem) (*string, error) {
	collection := getCollection(client)
	defer disconnect(client)

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
	collection := getCollection(client)
	defer disconnect(client)

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

func findById(c context.Context, id string) (*BlogItem, error) {
	collection := getCollection(client)
	defer disconnect(client)

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Could nit parse Id")
	}
	blogEntity := BlogItem{}
	err = collection.FindOne(c, bson.M{"_id": oid}).Decode(&blogEntity)
	if err == mongo.ErrNoDocuments {
		log.Printf("Not found %v\n", err)
		return nil, status.Errorf(codes.NotFound, "Not found document with id %v\n", id)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server Error")
	}
	return &blogEntity, nil

}
