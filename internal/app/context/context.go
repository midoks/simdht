package context

import (
	"fmt"

	"github.com/valyala/fasthttp"

	"github.com/midoks/simdht/internal/render"
)

var nowCtx *fasthttp.RequestCtx

func SetCtx(ctx *fasthttp.RequestCtx) {
	nowCtx = ctx
}

func HTML(status int, name string) {

	if nowCtx != nil {
		content, _ := render.HTML(name)
		code, err := nowCtx.Write(content)
		if err != nil {
			fmt.Println("HTML", code, err)
		}
	}
}

func Handler(routerHander fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		SetCtx(ctx)

		routerHander(ctx)
	}
}
