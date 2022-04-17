package context

import (
	"fmt"

	"github.com/valyala/fasthttp"

	"github.com/midoks/simdht/internal/render"
)

// Context represents context of a request.
type Context struct {
	*fasthttp.RequestCtx
}

func HTML(status int, name string) {
	fmt.Println(status, name)
	render.HTML(status, name)
}

func Handler(routerHander fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		fmt.Println("cc")

		// c := &Context{
		// 	Context: ctx,
		// }

		routerHander(ctx)
	}
}
