package dht

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	sdht "github.com/midoks/simdht/internal/shiyanhui/dht"
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

	downloader := sdht.NewWire(65536, 65536, 65536)
	go func() {
		// once we got the request result
		for resp := range downloader.Response() {

			metadata, err := sdht.Decode(resp.MetadataInfo)
			if err != nil {
				continue
			}
			info := metadata.(map[string]interface{})

			if _, ok := info["name"]; !ok {
				continue
			}

			fmt.Println("info", info)

			bt := bitTorrent{
				InfoHash: hex.EncodeToString(resp.InfoHash),
				Name:     info["name"].(string),
			}

			if v, ok := info["files"]; ok {
				files := v.([]interface{})
				bt.Files = make([]file, len(files))

				for i, item := range files {
					f := item.(map[string]interface{})
					bt.Files[i] = file{
						Path:   f["path"].([]interface{}),
						Length: f["length"].(int),
					}
				}
			} else if _, ok := info["length"]; ok {
				bt.Length = info["length"].(int)
			}

			data, err := json.Marshal(bt)
			if err == nil {
				fmt.Printf("%s\n\n", data)
			}

		}
	}()
	go downloader.Run()

	config := sdht.NewCrawlConfig()
	config.OnAnnouncePeer = func(infoHash, ip string, port int) {
		// request to download the metadata info
		fmt.Println("announce peer:", hex.EncodeToString([]byte(infoHash)), ip, port)
		downloader.Request([]byte(infoHash), ip, port)
	}

	config.CheckKBucketPeriod = time.Duration(time.Second * 3)
	d := sdht.New(config)

	go func() {
		for {
			fmt.Println(d.PrimeNodes)
			time.Sleep(time.Second * 3)
		}
	}()
	d.Run()
}
