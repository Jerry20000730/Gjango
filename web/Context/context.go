package context

import "net/http"

type Context struct {
	W http.ResponseWriter
	R *http.Request
}

func (c *Context) HTML(statusCode int, html string) error {
	// return a 200 ok
	c.W.WriteHeader(statusCode)

	// in the http response header - content-type: text/html; charset=utf-8
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")

	_, err := c.W.Write([]byte(html))

	return err
}
