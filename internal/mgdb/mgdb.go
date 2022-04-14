package mgdb

import (
	"context"
	// "fmt"
	// "time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	err error

	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
)

func Init() error {

	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	if client, err = mongo.Connect(context.TODO(), clientOptions); err != nil {
		return err
	}
	db = client.Database("simdht")
	collection = db.Collection("list")
	collection = collection

	return nil
}
