package dht

import (
	"bufio"
	"fmt"
	"net"

	sdht "github.com/shiyanhui/dht"
)

var connectResponse = []byte("HTTP/1.1 200 OK\r\n\r\n")

func handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		// read from the connection
		// ... ...
		// write to the connection
		//... ...
		reader := bufio.NewReader(conn)
		var buf [4096]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Println("read from client failed, err:", err)
			break
		}
		recvStr := string(buf[:n])

		fmt.Println("收到client端发来的数据：", recvStr)
		conn.Write(connectResponse)
		conn.Write([]byte("ddd24")) //发送数据
		break

	}
}

func Debug() {
	downloader := sdht.NewWire(65536, 1024, 256)
	go func() {
		// once we got the request result
		for resp := range downloader.Response() {
			fmt.Println("downloader", string(resp.InfoHash), string(resp.MetadataInfo))
		}
	}()
	go downloader.Run()

	config := sdht.NewCrawlConfig()
	config.OnAnnouncePeer = func(infoHash, ip string, port int) {
		// request to download the metadata info
		fmt.Println("an", infoHash, ip, port)
		downloader.Request([]byte(infoHash), ip, port)
	}
	d := sdht.New(config)

	d.Run()
}

func Run() {
	fmt.Println("DHT START")

	go Debug()

	l, err := net.Listen("tcp", ":8188")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			break
		}
		// start a new goroutine to handle
		// the new connection.
		go handleConn(c)
	}

}
