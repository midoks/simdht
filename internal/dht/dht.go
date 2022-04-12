package dht

import (
	"encoding/hex"
	"fmt"

	sdht "github.com/shiyanhui/dht"
)

var connectResponse = []byte("HTTP/1.1 200 OK\r\n\r\n")

type file struct {
	Path   []interface{} `json:"path"`
	Length int           `json:"length"`
}

type bitTorrent struct {
	InfoHash string `json:"infohash"`
	Name     string `json:"name"`
	Files    []file `json:"files,omitempty"`
	Length   int    `json:"length,omitempty"`
}

func Run() {
	fmt.Println("DHT START")

	downloader := sdht.NewWire(65536, 1024, 256)
	go func() {
		// once we got the request result
		for resp := range downloader.Response() {
			fmt.Println("downloader", hex.EncodeToString(resp.InfoHash))
			fmt.Println("downloader", string(resp.InfoHash), string(resp.MetadataInfo))
		}
	}()
	go downloader.Run()

	config := sdht.NewCrawlConfig()
	config.OnAnnouncePeer = func(infoHash, ip string, port int) {
		// request to download the metadata info
		fmt.Println(hex.EncodeToString([]byte(infoHash)))
		fmt.Println("announce peer:", infoHash, ip, port)
		downloader.Request([]byte(infoHash), ip, port)
	}
	d := sdht.New(config)

	go d.Run()

}
