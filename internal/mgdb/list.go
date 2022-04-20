package mgdb

import (
	"errors"
	"fmt"

	"github.com/qiniu/qmgo"
)

type File struct {
	Path   []interface{} `bson:"path";json:"path"`
	Length int           `bson:"length"`
}

type BitTorrent struct {
	InfoHash string `bson:"infohash"`
	Name     string `bson:"name"`
	Files    []File `bson:"files,omitempty"`
	Length   int    `bson:"length,omitempty"`
}

func AddTorrent(data BitTorrent) (result *qmgo.InsertOneResult, err error) {
	if collection != nil {
		result, err = collection.InsertOne(ctx, data)
		if err != nil {
			return nil, err
		}
		return result, err
	}

	return nil, errors.New("mongo disconnected!")
}

func DeleteTorrent() {
}

func Debug() {

	fmt.Println("ddd")
}
