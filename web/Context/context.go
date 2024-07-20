package context

import (
	"github.com/Jerry20000730/Gjango/web/Render"
	"net/http"
	"net/url"
)

type Context struct {
	W          http.ResponseWriter
	R          *http.Request
	queryCache url.Values
}

// initQueryCache initializes the query cache for the Context if it hasn't been initialized yet.
// This method checks if the queryCache field is nil, indicating that it hasn't been initialized.
// If it is nil, it then checks if there are any query parameters present in the request URL.
// If query parameters are present, it assigns them to the queryCache. Otherwise, it initializes
// queryCache as an empty url.Values object. This ensures that subsequent accesses to query parameters
// do not need to parse the URL again, improving performance for multiple query parameter accesses.
func (c *Context) initQueryCache() {
	if c.queryCache == nil {
		if c.R.URL.Query() != nil {
			c.queryCache = c.R.URL.Query()
		} else {
			c.queryCache = url.Values{}
		}
	}
}

// GetQuery retrieves the first value associated with the specified query parameter key.
// It ensures the query cache is initialized before attempting to retrieve the value.
// If the key exists in the query parameters, its first value is returned. If the key
// does not exist, an empty string is returned. This method is useful for query parameters
// that are expected to have a single value.
//
// Parameters:
//   - key: The query parameter key to look up.
//
// Returns:
//   - The first value associated with the key if it exists; otherwise, an empty string.
func (c *Context) GetQuery(key string) string {
	c.initQueryCache()
	return c.queryCache.Get(key)
}

// GetQueryArray retrieves the values associated with the specified query parameter key as a slice of strings.
// It also returns a boolean indicating whether the key was found in the query parameters.
//
// This method first ensures that the query cache is initialized by calling initQueryCache.
// It then attempts to retrieve the values associated with the provided key from the query cache.
// If the key exists, the method returns the slice of strings and true.
// If the key does not exist, it returns an empty slice of strings and false.
//
// Parameters:
//   - key: The query parameter key to look up.
//
// Returns:
//   - A slice of strings containing the values associated with the key.
//   - A boolean indicating whether the key was found in the query parameters.
func (c *Context) GetQueryArray(key string) ([]string, bool) {
	c.initQueryCache()
	values, ok := c.queryCache[key]
	return values, ok
}

// QueryArray retrieves all values associated with the specified key from the query parameters.
// It initializes the query cache if it's not already initialized, then looks up the key in the cache.
// If the key exists, it returns all associated values as a slice of strings. If the key does not exist,
// an empty slice is returned. This method is useful for query parameters that can have multiple values.
//
// Parameters:
//
//	key - The query parameter key to retrieve values for.
//
// Returns:
//
//	A slice of strings containing all values associated with the key. Returns an empty slice if the key does not exist.
func (c *Context) QueryArray(key string) []string {
	c.initQueryCache()
	values, _ := c.queryCache[key]
	return values
}

// GetQueryWithDefaultValues
func (c *Context) GetQueryWithDefaultValues(key string, defaultValue string) string {
	c.initQueryCache()
	values, ok := c.queryCache[key]
	if !ok {
		return defaultValue
	}
	return values[0]
}

// Render a general render function for rendering different types of data
// by passing a specific render as the second parameter.
// This method delegates the actual rendering process to the passed Render interface implementation,
// allowing for flexible rendering of various content types such as strings, JSON, XML, etc.
// The HTTP status code is set before rendering the content.
//
// Parameters:
//   - code: HTTP status code to be set for the response.
//   - r: An implementation of the Render interface to render the response.
//
// Returns:
//   - An error if the rendering process fails, otherwise nil.
func (c *Context) Render(code int, r Render.Render) error {
	err := r.Render(c.W)
	c.W.WriteHeader(code)
	return err
}

// String is a render function for rendering string content on the website.
// It utilizes the Render method to render string content, formatted according to the provided format string
// and values. This method is a convenience wrapper around the Render method for string rendering.
//
// Parameters:
//   - status: HTTP status code to be set for the response.
//   - format: A format string to be used for rendering the string content.
//   - values: A variadic parameter of values to be formatted according to the format string.
//
// Returns:
//   - An error if the rendering process fails, otherwise nil.
func (c *Context) String(status int, format string, values ...any) (err error) {
	err = c.Render(status, Render.StringRender{
		Format: format,
		Data:   values,
	})
	return
}

// XML is a render function for rendering XML data on the website.
// It delegates the rendering process to the Render method, specifically for XML content.
// This method is useful for APIs or services that need to return XML formatted responses.
//
// Parameters:
//   - status: HTTP status code to be set for the response.
//   - data: The data to be rendered as XML.
//
// Returns:
//   - An error if the rendering process fails, otherwise nil.
func (c *Context) XML(status int, data any) error {
	return c.Render(status, Render.XMLRender{Data: data})
}

// JSON is a render function for rendering JSON data on the website.
// It uses the Render method to render data in JSON format, suitable for REST APIs and web services
// that communicate using JSON.
//
// Parameters:
//   - status: HTTP status code to be set for the response.
//   - data: The data to be rendered as JSON.
//
// Returns:
//   - An error if the rendering process fails, otherwise nil.
func (c *Context) JSON(status int, data any) error {
	return c.Render(status, Render.JSONRender{Data: data})
}

// HTML is a render function for rendering HTML data on the website.
// This method is designed for rendering raw HTML content. It leverages the Render method,
// passing an HTMLRender instance configured for non-template HTML rendering.
//
// Parameters:
//   - status: HTTP status code to be set for the response.
//   - html: The raw HTML content to be rendered.
//
// Returns:
//   - An error if the rendering process fails, otherwise nil.
func (c *Context) HTML(status int, html string) error {
	return c.Render(status, Render.HTMLRender{IsTemplate: false, Data: html})
}

// HTMLTemplate is a render function for rendering HTML templates on the website.
// It provides a way to render dynamic HTML content using templates. This method is useful
// for web applications that serve HTML pages with dynamic data.
//
// Parameters:
//   - preloader: A preloader containing the HTML template to be rendered.
//   - status: HTTP status code to be set for the response.
//   - name: The name of the template to be rendered.
//   - data: The data to be passed to the template.
//
// Returns:
//   - An error if the rendering process fails, otherwise nil.
func (c *Context) HTMLTemplate(preloader *Render.HTMLPreloader, status int, name string, data any) error {
	return c.Render(status, Render.HTMLRender{
		IsTemplate: true,
		Name:       name,
		Data:       data,
		Template:   preloader.Template,
	})
}

// Redirect is a method for redirecting to another URL.
// Technically not a rendering function, this method uses the Render method to perform an HTTP redirect.
// It is useful for web applications that need to redirect users to different URLs based on certain conditions.
//
// Parameters:
//   - status: HTTP status code to be used for the redirect, typically 302 for temporary redirects.
//   - location: The URL to redirect to.
//
// Returns:
//   - An error if the redirect process fails, otherwise nil.
func (c *Context) Redirect(status int, location string) error {
	return c.Render(status, Render.Redirect{
		Code:     status,
		Request:  c.R,
		Location: location,
	})
}
