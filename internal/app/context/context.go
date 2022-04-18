package context

import (
	"fmt"

	"github.com/valyala/fasthttp"

	"github.com/midoks/simdht/internal/render"
)

func Base(ctx *fasthttp.RequestCtx) {

}

func HTML(ctx *fasthttp.RequestCtx, status int, name string, data interface{}) {

	content, _ := render.HTML(name, data)
	code, err := ctx.Write(content)
	if err != nil {
		fmt.Println("HTML", code, err)
	}

}

func Handler(routerHander fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		routerHander(ctx)
	}
}
