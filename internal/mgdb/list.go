package mgdb

import (
	"fmt"

	"github.com/qiniu/qmgo"
)

type File struct {
	Path   []interface{} `json:"path"`
	Length int           `json:"length"`
}

type BitTorrent struct {
	InfoHash string `json:"infohash"`
	Name     string `json:"name"`
	Files    []File `json:"files,omitempty"`
	Length   int    `json:"length,omitempty"`
}

func AddTorrent(data BitTorrent) (result *qmgo.InsertOneResult, err error) {
	return collection.InsertOne(ctx, data)
}

func DeleteTorrent() {
}

func Debug() {

	fmt.Println("ddd")
}
