package cmd

import (
	"fmt"
	"log"

	"github.com/fasthttp/router"
	"github.com/urfave/cli"
	"github.com/valyala/fasthttp"
)

var Service = cli.Command{
	Name:        "service",
	Usage:       "This command starts all services",
	Description: `Start DHT services`,
	Action:      runAllService,
	Flags: []cli.Flag{
		stringFlag("config, c", "", "Custom configuration file path"),
	},
}

func runAllService(c *cli.Context) error {
	// err := router.GlobalInit(c.String("config"))
	// if err != nil {
	// 	log.Errorf("Failed to initialize application: %s", err)
	// }

	// app.Start(conf.Web.HttpPort)

	r := router.New()
	r.GET("/", Index)
	r.GET("/hello/{name}", Hello)

	if err := fasthttp.ListenAndServe(":11010", r.Handler); err != nil {
		log.Fatalf("Error in ListenAndServe: %v", err)
	}
	return nil
}

func Hello(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, %s!\n", ctx.UserValue("name"))
}
func Index(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, world!\n\n")

	fmt.Fprintf(ctx, "Request method is %q\n", ctx.Method())
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
}
