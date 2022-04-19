package context

import (
	// "fmt"

	"github.com/valyala/fasthttp"

	"github.com/midoks/simdht/internal/render"
)

const (
	_CONTENT_TYPE    = "Content-Type"
	_CONTENT_BINARY  = "application/octet-stream"
	_CONTENT_JSON    = "application/json"
	_CONTENT_HTML    = "text/html"
	_CONTENT_PLAIN   = "text/plain"
	_CONTENT_XHTML   = "application/xhtml+xml"
	_CONTENT_XML     = "text/xml"
	_DEFAULT_CHARSET = "UTF-8"
)

func Base(ctx *fasthttp.RequestCtx) {

}

func HTML(ctx *fasthttp.RequestCtx, status int, name string, data interface{}) {

	content, _ := render.HTML(name, data)

	ctx.Response.Header.SetStatusCode(status)
	ctx.Response.Header.SetContentType(_CONTENT_HTML)
	ctx.Write(content)
}

func Handler(routerHander fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		routerHander(ctx)
	}
}
