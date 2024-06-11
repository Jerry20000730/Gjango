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
	g.Get("/hello", func(ctx *context.Context) {
		fmt.Fprintf(ctx.W, "<h1>Welcome to Gjango</h1> <p>You have successfully initiate the Gjango web framework</p>")
	})
	g.Post("/hello", func(ctx *context.Context) {
		fmt.Fprintf(ctx.W, "<h1>Welcome to Gjango</h1> <p>This is a POST request and you have successfully initiate the Gjango web framework</p>")
	})
	g.Get("/get/:id", func(ctx *context.Context) {
		fmt.Fprintf(ctx.W, "<h1>Welcome to Gjango</h1> <p>This is a GET request and you get user info path variable")
	})
	engine.Run()
}
