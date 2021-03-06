package app

import (
	"fmt"
	// "net/http"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"github.com/midoks/simdht/internal/app/context"
	"github.com/midoks/simdht/internal/app/router/admin"
)

func Start(port string) {

	r := router.New()
	r.GET("/", Index)
	r.GET("/s/{KW}.html", Soso)
	r.GET("/s/{KW}/{Cat}/{Sort}.html", Soso)
	r.GET("/s/{KW}/{Cat}/{Sort}/{Pos}.html", Soso)
	r.GET("/info/{HASH_ID}", Info)
	r.GET("/gc", admin.GcInfo)
	r.GET("/debug", Debug)

	r.NotFound = fasthttp.FSHandler("./public", 0)

	if err := fasthttp.ListenAndServe(":"+port, context.Handler(r.Handler)); err != nil {
		fmt.Errorf("Error in fasthttp.ListenAndServe: %v", err)
	}
}

func Index(ctx *fasthttp.RequestCtx) {
	_this := context.Base(ctx)
	_this.HTML(200, "index")
}

func Soso(ctx *fasthttp.RequestCtx) {
	_this := context.Base(ctx)

	keyword := ctx.UserValue("KW")
	cat := ctx.UserValue("Cat")
	sort := ctx.UserValue("Sort")

	fmt.Println("keyword:", keyword)
	fmt.Println("cat:", cat)
	fmt.Println("sort:", sort)

	_this.Data["Kw"] = keyword
	_this.Data["Name"] = "ww"

	_this.HTML(200, "search/index")
}

func Info(ctx *fasthttp.RequestCtx) {
	_this := context.Base(ctx)
	_this.HTML(200, "search/info")
}

func Hello(ctx *fasthttp.RequestCtx) {
	_this := context.Base(ctx)
	_this.HTML(200, "info")
}

func Debug(ctx *fasthttp.RequestCtx) {

	newCtx := &fasthttp.Request{}
	ctx.Request.CopyTo(newCtx)

	fmt.Fprintf(ctx, "Hello, world!\n\n")
	fmt.Fprintf(ctx, "Request method is %s\n", ctx.Method())
	fmt.Fprintf(ctx, "RequestURI is %q\n", ctx.RequestURI())
	fmt.Fprintf(ctx, "Requested path is %q\n", ctx.Path())
	fmt.Fprintf(ctx, "Host is %q\n", ctx.Host())
	fmt.Fprintf(ctx, "Query string is %q\n", ctx.QueryArgs())
	fmt.Fprintf(ctx, "User-Agent is %q\n", ctx.UserAgent())
	fmt.Fprintf(ctx, "Connection has been established at %s\n", ctx.ConnTime())
	fmt.Fprintf(ctx, "Request has been started at %s\n", ctx.Time())
	fmt.Fprintf(ctx, "Serial request number for the current connection is %d\n", ctx.ConnRequestNum())
	fmt.Fprintf(ctx, "Your ip is %q\n\n", ctx.RemoteIP())

	fmt.Fprintf(ctx, "Raw request is:\n---CUT---\n%s\n---CUT---", &ctx.Request)

	ctx.SetContentType("text/plain; charset=utf8")

	// Set arbitrary headers
	ctx.Response.Header.Set("X-My-Header", "my-header-value")

	// Set cookies
	var c fasthttp.Cookie
	c.SetKey("cookie-name")
	c.SetValue("cookie-value")
	ctx.Response.Header.SetCookie(&c)
	newCtx = nil

	admin.GcInfo(ctx)
}
