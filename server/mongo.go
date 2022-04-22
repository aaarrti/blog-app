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
	"os"
)

var client = func() *mongo.Client {
	mongoUri, ok := os.LookupEnv("MONGO_URI")
	if !ok {
		log.Fatalf("MONGO_URI env var not set")
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatal(err)
	}
	return client
}()

func testConnection() {
	err := client.Connect(context.TODO())
	if err != nil {
		log.Fatalf("Failed to connect to mongo db: %v\n", err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Error while pinging mongo db: %v\n", err)
	}
	log.Printf("Connected successfully to MongoDB")
}

func getCollection(client *mongo.Client) *mongo.Collection {
	_collection := client.Database("blog-db").Collection("blogs")
	return _collection
}

func save(ctx context.Context, item *CreateBlogItem) (*string, error) {
	collection := getCollection(client)

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

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Could not parse Id")
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

func deleteById(c context.Context, id string) error {
	collection := getCollection(client)

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Could not parse Id")
	}
	res, err := collection.DeleteOne(c, bson.M{"_id": oid})
	if err != nil {
		log.Printf("Error while deleting %v\n", err)
		return status.Errorf(codes.Internal, "Error while deleting")
	}
	if res.DeletedCount == 0 {
		return status.Errorf(codes.NotFound, "Not found document with id %v\n", id)
	}
	return nil
}

func findAll() ([]BlogItem, error) {
	collection := getCollection(client)
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Printf("Error while getting all blogs")
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}
	var res []BlogItem
	for cursor.Next(context.TODO()) {
		var document BlogItem
		err = cursor.Decode(&document)
		log.Print(document)
		res = append(res, document)
	}
	return res, nil
}
