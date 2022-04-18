package app

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/midoks/simdht/internal/app/router/admin"
	"github.com/midoks/simdht/internal/app/template"
	"github.com/midoks/simdht/internal/render"
)

var defaultClientsCount = runtime.NumCPU()

var (
	fakeResponse = []byte("Hello, world!")
	getRequest   = "GET /foobar?baz HTTP/1.1\r\nHost: google.com\r\nUser-Agent: aaa/bbb/ccc/ddd/eee Firefox Chrome MSIE Opera\r\n" +
		"Referer: http://example.com/aaa?bbb=ccc\r\nCookie: foo=bar; baz=baraz; aa=aakslsdweriwereowriewroire\r\n\r\n"
	postRequest = fmt.Sprintf("POST /foobar?baz HTTP/1.1\r\nHost: google.com\r\nContent-Type: foo/bar\r\nContent-Length: %d\r\n"+
		"User-Agent: Opera Chrome MSIE Firefox and other/1.2.34\r\nReferer: http://google.com/aaaa/bbb/ccc\r\n"+
		"Cookie: foo=bar; baz=baraz; aa=aakslsdweriwereowriewroire\r\n\r\n%q",
		len(fakeResponse), fakeResponse)
)

type RequestHandler func(ctx *fasthttp.RequestCtx)

type realServer interface {
	Serve(ln net.Listener) error
}

type fakeServerConn struct {
	net.TCPConn
	ln            *fakeListener
	requestsCount int
	pos           int
	closed        uint32
}

func (c *fakeServerConn) Read(b []byte) (int, error) {
	nn := 0
	reqLen := len(c.ln.request)
	for len(b) > 0 {
		if c.requestsCount == 0 {
			if nn == 0 {
				return 0, io.EOF
			}
			return nn, nil
		}
		pos := c.pos % reqLen
		n := copy(b, c.ln.request[pos:])
		b = b[n:]
		nn += n
		c.pos += n
		if n+pos == reqLen {
			c.requestsCount--
		}
	}
	return nn, nil
}

func (c *fakeServerConn) Write(b []byte) (int, error) {
	return len(b), nil
}

var fakeAddr = net.TCPAddr{
	IP:   []byte{1, 2, 3, 4},
	Port: 12345,
}

func (c *fakeServerConn) RemoteAddr() net.Addr {
	return &fakeAddr
}

func (c *fakeServerConn) Close() error {
	if atomic.AddUint32(&c.closed, 1) == 1 {
		c.ln.ch <- c
	}
	return nil
}

func (c *fakeServerConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *fakeServerConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func initTmplConf() {
	path, _ := os.Getwd()
	_path := filepath.Dir(path)
	_path = filepath.Dir(_path)

	renderOpt := render.Options{
		Directory:         filepath.Join(_path, "templates"),
		AppendDirectories: []string{filepath.Join(_path+"/custom", "templates")},
		Funcs:             template.FuncMap(),
		IndentJSON:        true,
		Delims: render.Delims{
			Left:  "{{",
			Right: "}}",
		},
	}

	render.Renderer(renderOpt)
}

// go test -bench=. -benchmem -benchtime=1s

// go test -bench=BenchmarkServerGet -benchmem -benchtime=1s
func BenchmarkServerGetGcInfo(b *testing.B) {
	benchmarkServerGetCallback(b, defaultClientsCount, 1, func(ctx *fasthttp.RequestCtx) {
		admin.GcInfo(ctx)
	})
}

func BenchmarkServerGetIndex(b *testing.B) {
	benchmarkServerGetCallback(b, defaultClientsCount, 1, func(ctx *fasthttp.RequestCtx) {
		Index(ctx)
	})
}

func BenchmarkServerGetHello(b *testing.B) {

	initTmplConf()

	benchmarkServerGetCallback(b, defaultClientsCount, 1, func(ctx *fasthttp.RequestCtx) {
		Hello(ctx)
	})
}

func benchmarkServerGetCallback(b *testing.B, clientsCount, requestsPerConn int, handler RequestHandler) {
	ch := make(chan struct{}, b.N)
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			if !ctx.IsGet() {
				b.Fatalf("Unexpected request method: %q", ctx.Method())
			}
			handler(ctx)
			if requestsPerConn == 1 {
				ctx.SetConnectionClose()
			}
			registerServedRequest(b, ch)
		},
		Concurrency: 16 * clientsCount,
	}
	benchmarkServer(b, s, clientsCount, requestsPerConn, getRequest)
	verifyRequestsServed(b, ch)
}

func benchmarkServerGet(b *testing.B, clientsCount, requestsPerConn int) {
	ch := make(chan struct{}, b.N)
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			if !ctx.IsGet() {
				b.Fatalf("Unexpected request method: %q", ctx.Method())
			}

			admin.GcInfo(ctx)

			if requestsPerConn == 1 {
				ctx.SetConnectionClose()
			}
			registerServedRequest(b, ch)
		},
		Concurrency: 16 * clientsCount,
	}
	benchmarkServer(b, s, clientsCount, requestsPerConn, getRequest)
	verifyRequestsServed(b, ch)
}

type fakeListener struct {
	lock            sync.Mutex
	requestsCount   int
	requestsPerConn int
	request         []byte
	ch              chan *fakeServerConn
	done            chan struct{}
	closed          bool
}

func (ln *fakeListener) Accept() (net.Conn, error) {
	ln.lock.Lock()
	if ln.requestsCount == 0 {
		ln.lock.Unlock()
		for len(ln.ch) < cap(ln.ch) {
			time.Sleep(10 * time.Millisecond)
		}
		ln.lock.Lock()
		if !ln.closed {
			close(ln.done)
			ln.closed = true
		}
		ln.lock.Unlock()
		return nil, io.EOF
	}
	requestsCount := ln.requestsPerConn
	if requestsCount > ln.requestsCount {
		requestsCount = ln.requestsCount
	}
	ln.requestsCount -= requestsCount
	ln.lock.Unlock()

	c := <-ln.ch
	c.requestsCount = requestsCount
	c.closed = 0
	c.pos = 0

	return c, nil
}

func (ln *fakeListener) Close() error {
	return nil
}

func (ln *fakeListener) Addr() net.Addr {
	return &fakeAddr
}

func newFakeListener(requestsCount, clientsCount, requestsPerConn int, request string) *fakeListener {
	ln := &fakeListener{
		requestsCount:   requestsCount,
		requestsPerConn: requestsPerConn,
		request:         []byte(request),
		ch:              make(chan *fakeServerConn, clientsCount),
		done:            make(chan struct{}),
	}
	for i := 0; i < clientsCount; i++ {
		ln.ch <- &fakeServerConn{
			ln: ln,
		}
	}
	return ln
}

func registerServedRequest(b *testing.B, ch chan<- struct{}) {
	select {
	case ch <- struct{}{}:
	default:
		b.Fatalf("More than %d requests served", cap(ch))
	}
}

func benchmarkServer(b *testing.B, s realServer, clientsCount, requestsPerConn int, request string) {
	ln := newFakeListener(b.N, clientsCount, requestsPerConn, request)
	ch := make(chan struct{})
	go func() {
		s.Serve(ln) //nolint:errcheck
		ch <- struct{}{}
	}()

	<-ln.done

	select {
	case <-ch:
	case <-time.After(10 * time.Second):
		b.Fatalf("Server.Serve() didn't stop")
	}
}

func verifyRequestsServed(b *testing.B, ch <-chan struct{}) {
	requestsServed := 0
	for len(ch) > 0 {
		<-ch
		requestsServed++
	}
	requestsSent := b.N
	for requestsServed < requestsSent {
		select {
		case <-ch:
			requestsServed++
		case <-time.After(100 * time.Millisecond):
			b.Fatalf("Unexpected number of requests served %d. Expected %d", requestsServed, requestsSent)
		}
	}
}
