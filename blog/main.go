package main

import (
	"fmt"
	"github.com/Jerry20000730/Gjango/web"
	context "github.com/Jerry20000730/Gjango/web/Context"
)

func main() {
	engine := web.NewEngine()
	fmt.Println("[INFO] Gjango is listening on port: " + engine.GetPort())
	g := engine.Router.NewGroup("user")
	g.MiddlewareBind(func(next web.Handler) web.Handler {
		return func(ctx *context.Context) {
			// pre-middleware
			fmt.Println("Pre Handler")
			next(ctx)
			// post-middleware
			fmt.Println("Post Handler")
		}
	})
	g.Get("/hello/get", func(ctx *context.Context) {
		_, _ = fmt.Fprintf(ctx.W, "<h1>Welcome to Gjango</h1> <p>You have successfully initiate the Gjango web framework</p>")
	})
	g.Post("/hello", func(ctx *context.Context) {
		_, _ = fmt.Fprintf(ctx.W, "<h1>Welcome to Gjango</h1> <p>This is a POST request and you have successfully initiate the Gjango web framework</p>")
	})
	g.Get("/get/:id", func(ctx *context.Context) {
		_, _ = fmt.Fprintf(ctx.W, "<h1>Welcome to Gjango</h1> <p>This is a GET request and you get user info path variable")
	})
	g.Get("/hello/*/get", func(ctx *context.Context) {
		_, _ = fmt.Fprintf(ctx.W, "<h1>Welcome to Gjango</h1> <p>This is a GET request and I don't know what you are looking for</p>")
	})
	engine.Run()
}
