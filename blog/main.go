package main

import (
	"fmt"
	"github.com/Jerry20000730/Gjango/web"
	context "github.com/Jerry20000730/Gjango/web/Context"
	"log"
	"net/http"
)

func BlogLog(next web.Handler) web.Handler {
	return func(ctx *context.Context) {
		fmt.Println("BlogLog begin for get")
		next(ctx)
		fmt.Println("BlogLog end")
	}
}

type User struct {
	Name      string   `xml:"name" json:"name"`
	Age       int      `xml:"age" json:"age"`
	Addresses []string `xml:"addresses" json:"addresses"`
	Email     string   `xml:"email" json:"email" gjango:"required"`
}

func main() {
	engine := web.NewEngine()
	fmt.Println("[INFO] Gjango is listening on port: " + engine.GetPort())
	g := engine.Router.NewGroup("user")
	g.MiddlewareRegister(func(next web.Handler) web.Handler {
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
	}, BlogLog)
	g.Post("/hello", func(ctx *context.Context) {
		_, _ = fmt.Fprintf(ctx.W, "<h1>Welcome to Gjango</h1> <p>This is a POST request and you have successfully initiate the Gjango web framework</p>")
	})
	g.Get("/get/:id", func(ctx *context.Context) {
		_, _ = fmt.Fprintf(ctx.W, "<h1>Welcome to Gjango</h1> <p>This is a GET request and you get user info path variable")
	})
	g.Get("/hello/*/get", func(ctx *context.Context) {
		_, _ = fmt.Fprintf(ctx.W, "<h1>Welcome to Gjango</h1> <p>This is a GET request and I don't know what you are looking for</p>")
	})
	g.Get("/get/html", func(ctx *context.Context) {
		_ = ctx.HTML(http.StatusOK, "<h1>HTML Template</h1><p>This is a template for html and test if the html is successfully returned and rendered</p>")
	})

	engine.PreLoadTemplate("template/*.html")
	g.Get("/template", func(ctx *context.Context) {
		user := &User{
			Name: "jerry",
		}
		err := ctx.HTMLTemplate(&engine.HTMLPreloader, http.StatusOK, "login.html", user)
		if err != nil {
			log.Println(err)
		}
	})

	g.Get("/json", func(ctx *context.Context) {
		user := &User{
			Name: "jerry",
		}
		err := ctx.JSON(http.StatusOK, user)
		if err != nil {
			log.Println(err)
		}
	})

	g.Get("/xml", func(ctx *context.Context) {
		user := &User{
			Name: "jerry",
		}
		err := ctx.XML(http.StatusOK, user)
		if err != nil {
			log.Println(err)
		}
	})

	g.Get("/download", func(ctx *context.Context) {
		engine.FileManager.FileAttachment(ctx, "template/test.xlsx", "xxxx.xlsx")
	})
	g.Get("/fs", func(ctx *context.Context) {
		engine.FileManager.FileFromFileSystem(ctx, "test.xlsx", http.Dir("template"))
	})
	g.Get("/redirect", func(ctx *context.Context) {
		// status must be 302
		err := ctx.Redirect(http.StatusFound, "/user/template")
		if err != nil {
			log.Println(err)
		}
	})
	g.Get("/string", func(ctx *context.Context) {
		_ = ctx.String(http.StatusOK, "Test %s gjango web framework, int can also be passed: %d", "self-designed", 1)
	})
	g.Get("/queryMap", func(ctx *context.Context) {
		queryMap, _ := ctx.GetQueryMap("user")
		ctx.JSON(http.StatusOK, queryMap)
	})
	g.Post("/formPost", func(ctx *context.Context) {
		m, _ := ctx.GetPostFormArray("name")
		ctx.JSON(http.StatusOK, m)
	})
	g.Post("/formPostMap", func(ctx *context.Context) {
		m, _ := ctx.GetPostFormMap("user")
		ctx.JSON(http.StatusOK, m)
	})
	g.Post("/file", func(ctx *context.Context) {
		file, err := ctx.FormFile("file")
		if err != nil {
			log.Println(err)
		}
		err = ctx.SaveUploadedFile(file, "./upload/"+file.Filename)
		if err != nil {
			log.Println(err)
		}
	})
	g.Post("/multiFile", func(ctx *context.Context) {
		m, _ := ctx.GetPostFormMap("user")
		headers := ctx.FormFiles("file")
		for _, file := range headers {
			ctx.SaveUploadedFile(file, "./upload/"+file.Filename)
		}
		ctx.JSON(http.StatusOK, m)
	})
	g.Post("/jsonParse", func(ctx *context.Context) {
		user := &User{}
		err := ctx.ParseJSON(user, true, true)
		if err != nil {
			log.Println(err)
		}
		ctx.JSON(http.StatusOK, user)
	})
	engine.Run()
}
