package mgdb

import (
	"context"

	"github.com/qiniu/qmgo"
)

var (
	err        error
	ctx        context.Context
	client     *qmgo.Client
	db         *qmgo.Database
	collection *qmgo.Collection
)

func Init() error {

	ctx = context.Background()
	client, err = qmgo.NewClient(ctx, &qmgo.Config{Uri: "mongodb://127.0.0.1:27017"})

	if err != nil {
		return err
	}

	db = client.Database("simdht")
	collection = db.Collection("list")

	return nil
}
