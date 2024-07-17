package context

import (
	"github.com/Jerry20000730/Gjango/web/Render"
	"net/http"
)

type Context struct {
	W http.ResponseWriter
	R *http.Request
}

func (c *Context) Render(code int, r Render.Render) error {
	err := r.Render(c.W)
	c.W.WriteHeader(code)
	return err
}

func (c *Context) String(status int, format string, values ...any) (err error) {
	err = c.Render(status, Render.StringRender{
		Format: format,
		Data:   values,
	})
	return
}

func (c *Context) XML(status int, data any) error {
	return c.Render(status, Render.XMLRender{Data: data})
}

func (c *Context) JSON(status int, data any) error {
	return c.Render(status, Render.JSONRender{Data: data})
}

func (c *Context) HTML(status int, html string) error {
	return c.Render(status, Render.HTMLRender{IsTemplate: false, Data: html})
}

func (c *Context) HTMLTemplate(preloader *Render.HTMLPreloader, status int, name string, data any) error {
	return c.Render(status, Render.HTMLRender{
		IsTemplate: true,
		Name:       name,
		Data:       data,
		Template:   preloader.Template,
	})
}

func (c *Context) Redirect(status int, location string) error {
	return c.Render(status, Render.Redirect{
		Code:     status,
		Request:  c.R,
		Location: location,
	})
}
