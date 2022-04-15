package context

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

// Context represents context of a request.
type Context struct {
	*fasthttp.RequestCtx
}

func (*Context) HTML(status int, name string) {
	fmt.Println(status, name)
}
