package context

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jerry20000730/Gjango/web/Constant"
	"github.com/Jerry20000730/Gjango/web/Render"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
)

type Context struct {
	W          http.ResponseWriter
	R          *http.Request
	queryCache url.Values
	formCache  url.Values
}

// initQueryCache initializes the query cache for the Context if it hasn't been initialized yet.
// If query parameters are present, it assigns them to the queryCache. Otherwise, it initializes
// queryCache as an empty url.Values object. This ensures that subsequent accesses to query parameters
// do not need to parse the URL again, improving performance for multiple query parameter accesses.
func (c *Context) initQueryCache() {
	if c.R.URL.Query() != nil {
		c.queryCache = c.R.URL.Query()
	} else {
		c.queryCache = url.Values{}
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

// GetQueryWithDefaultValues similar to GetQuery, but with default value initiated
func (c *Context) GetQueryWithDefaultValues(key string, defaultValue string) string {
	c.initQueryCache()
	values, ok := c.queryCache[key]
	if !ok {
		return defaultValue
	}
	return values[0]
}

// QueryMap retrieves a map of string keys to string values for a given query parameter key.
// This method is useful for parsing query parameters with nested keys, such as "user[id]=1&user[age]=29".
// It leverages the GetQueryMap method to initialize the query cache (if not already done) and
// to validate and extract the query parameters into a map.
//
// Parameters:
//   - key: The base query parameter key to look for nested keys in.
//
// Returns:
//   - A map of string keys to string values representing the nested query parameters.
func (c *Context) QueryMap(key string) map[string]string {
	dicts, _ := c.GetQueryMap(key)
	return dicts
}

// GetQueryMap initializes the query cache (if necessary) and retrieves a map of string keys to string values
// for a given query parameter key, along with a boolean indicating whether the key was found.
// This method is specifically designed to handle query parameters with nested keys, parsing them into a map.
//
// Parameters:
//   - key: The base query parameter key to look for nested keys in.
//
// Returns:
//   - A map of string keys to string values representing the nested query parameters.
//   - A boolean indicating whether the key was found in the query parameters.
func (c *Context) GetQueryMap(key string) (map[string]string, bool) {
	c.initQueryCache()
	return c.getAndValidate(c.queryCache, key)
}

// getAndValidate parses a map of query parameters (map[string][]string) for nested keys based on a given key.
// It constructs a map of the nested keys to their first corresponding value. This method supports parsing
// query parameters formatted like "user[id]=1&user[age]=29", extracting the nested keys and their values.
//
// Parameters:
//   - m: The map of query parameters to parse.
//   - key: The base query parameter key to look for nested keys in.
//
// Returns:
//   - A map of string keys to string values representing the nested query parameters.
//   - A boolean indicating whether any nested keys were found for the given base key.
//
// Example:
//
//	For a query string "?user[id]=1&user[age]=29", calling getAndValidate with key "user"
//	would return a map {"id": "1", "age": "29"} and true.
func (c *Context) getAndValidate(m map[string][]string, key string) (map[string]string, bool) {
	dicts := make(map[string]string)
	exist := false
	for k, value := range m {
		if i := strings.IndexByte(k, '['); i >= 1 && k[0:i] == key {
			if j := strings.IndexByte(k[i+1:], ']'); j >= 1 {
				exist = true
				dicts[k[i+1:][:j]] = value[0]
			}
		}
	}
	return dicts, exist
}

// initPostFormCache initializes the form cache for the Context if it hasn't been initialized yet.
// This method parses the request body as multipart/form-data if the request's Content-Type
// indicates so and stores the parsed data in the formCache. If the request body is not
// multipart/form-data or if the request is nil, it initializes formCache as an empty url.Values object.
// This ensures that subsequent accesses to form parameters do not need to parse the request body again,
// improving performance for multiple form parameter accesses.
func (c *Context) initPostFormCache() {
	if c.R != nil {
		if err := c.R.ParseMultipartForm(Constant.DEFAULT_MAX_MEMORY); err != nil {
			if errors.Is(err, http.ErrNotMultipart) {
				log.Println(err)
			}
		}
		c.formCache = c.R.PostForm
	} else {
		c.formCache = url.Values{}
	}
}

// GetPostForm retrieves the first value associated with the specified form parameter key.
// It ensures the form cache is initialized before attempting to retrieve the value.
// If the key exists in the form parameters, its first value is returned along with a boolean true.
// If the key does not exist, an empty string and false are returned. This method is useful for form parameters
// that are expected to have a single value.
//
// Parameters:
//   - key: The form parameter key to look up.
//
// Returns:
//   - The first value associated with the key if it exists; otherwise, an empty string.
//   - A boolean indicating whether the key was found in the form parameters.
func (c *Context) GetPostForm(key string) (string, bool) {
	if values, ok := c.GetPostFormArray(key); ok {
		return values[0], ok
	}
	return "", false
}

// GetPostFormArray retrieves the values associated with the specified form parameter key as a slice of strings.
// It also returns a boolean indicating whether the key was found in the form parameters.
//
// This method first ensures that the form cache is initialized by calling initPostFormCache.
// It then attempts to retrieve the values associated with the provided key from the form cache.
// If the key exists, the method returns the slice of strings and true.
// If the key does not exist, it returns an empty slice of strings and false.
//
// Parameters:
//   - key: The form parameter key to look up.
//
// Returns:
//   - A slice of strings containing the values associated with the key.
//   - A boolean indicating whether the key was found in the form parameters.
func (c *Context) GetPostFormArray(key string) ([]string, bool) {
	c.initPostFormCache()
	values, ok := c.formCache[key]
	return values, ok
}

// PostFormArray retrieves all values associated with the specified key from the form parameters.
// It initializes the form cache if it's not already initialized, then looks up the key in the cache.
// If the key exists, it returns all associated values as a slice of strings. If the key does not exist,
// an empty slice is returned. This method is useful for form parameters that can have multiple values.
//
// Parameters:
//   - key: The form parameter key to retrieve values for.
//
// Returns:
//   - A slice of strings containing all values associated with the key. Returns an empty slice if the key does not exist.
func (c *Context) PostFormArray(key string) []string {
	values, _ := c.GetPostFormArray(key)
	return values
}

// GetPostFormMap initializes the form cache (if necessary) and retrieves a map of string keys to string values
// for a given form parameter key, along with a boolean indicating whether the key was found.
// This method is specifically designed to handle form parameters with nested keys, parsing them into a map.
//
// Parameters:
//   - key: The base form parameter key to look for nested keys in.
//
// Returns:
//   - A map of string keys to string values representing the nested form parameters.
//   - A boolean indicating whether the key was found in the form parameters.
func (c *Context) GetPostFormMap(key string) (map[string]string, bool) {
	c.initPostFormCache()
	return c.getAndValidate(c.formCache, key)
}

// FormFile retrieves a single file from the multipart form data based on the given name.
// It returns the file header if the file is found, allowing further operations like opening or reading the file.
// If the file is not found or an error occurs during retrieval, an error is returned.
//
// Parameters:
//   - name: The name attribute of the file input field from the form.
//
// Returns:
//   - *multipart.FileHeader: The file header for the specified file, if found.
//   - error: An error object if the file cannot be retrieved.
func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
	file, header, err := c.R.FormFile(name)
	if err != nil {
		return nil, err
	}
	err = file.Close()
	if err != nil {
		return nil, err
	}
	return header, nil
}

// FormFiles retrieves all files associated with the given name from the multipart form data.
// It returns a slice of file headers, each representing a file uploaded under the same name attribute in the form.
// If no files are found or an error occurs, an empty slice is returned.
//
// Parameters:
//   - name: The name attribute of the file input field from the form.
//
// Returns:
//   - []*multipart.FileHeader: A slice of file headers for the files associated with the specified name.
func (c *Context) FormFiles(name string) []*multipart.FileHeader {
	multipartForm, err := c.MultipartForm()
	if err != nil {
		return make([]*multipart.FileHeader, 0)
	}
	return multipartForm.File[name]
}

// MultipartForm parses the request body as multipart/form-data and returns the parsed multipart form.
// It ensures that the form data is parsed only once and caches the result for subsequent accesses.
// If parsing fails, an error is returned.
//
// Returns:
//   - *multipart.Form: The parsed multipart form containing all file and non-file form fields.
//   - error: An error object if parsing fails.
func (c *Context) MultipartForm() (*multipart.Form, error) {
	err := c.R.ParseMultipartForm(Constant.DEFAULT_MAX_MEMORY)
	return c.R.MultipartForm, err
}

// SaveUploadedFile saves an uploaded file to a specified destination on the server's filesystem.
// It opens the file associated with the provided file header, creates a new file at the destination path,
// and copies the contents of the uploaded file to the new file.
// If any step fails, an error is returned.
//
// Parameters:
//   - file: The file header of the uploaded file to be saved.
//   - dst: The destination path where the uploaded file should be saved.
//
// Returns:
//   - error: An error object if saving the file fails.
func (c *Context) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, src)
	return err
}

func (c *Context) ParseJSON(obj any, disallowUnknownField bool, isValidate bool) error {
	body := c.R.Body
	if body == nil {
		return errors.New("[ERROR] body is nil, invalid request")
	}
	decoder := json.NewDecoder(body)
	// if there is unknown fields
	// there will be errors
	if disallowUnknownField {
		decoder.DisallowUnknownFields()
	}
	if isValidate {
		err := validateJsonParam(obj, decoder)
		if err != nil {
			return err
		}
	} else {
		err := decoder.Decode(obj)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateJsonParam(obj any, decoder *json.Decoder) error {
	// parse to map, and then compare the key with the map key
	valueOf := reflect.ValueOf(obj)
	if valueOf.Kind() != reflect.Pointer {
		return errors.New("[ERROR] passing parameter is not a pointer")
	}
	elem := valueOf.Elem().Interface()
	of := reflect.ValueOf(elem)

	switch of.Kind() {
	case reflect.Struct:
		err := checkStructParams(obj, decoder, of)
		if err != nil {
			return err
		}
	case reflect.Slice, reflect.Array:
		elem := of.Type().Elem()
		elemType := elem.Kind()
		if elemType == reflect.Struct {
			return checkSliceParams(obj, elem, decoder)
		}
	default:
		err := decoder.Decode(obj)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkSliceParams(obj any, elem reflect.Type, decoder *json.Decoder) error {
	mapData := make([]map[string]interface{}, 0)
	err := decoder.Decode(&mapData)
	if err != nil {
		return err
	}
	if len(mapData) <= 0 {
		return nil
	}
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		gjangoTag := field.Tag.Get("gjango")
		tag := field.Tag.Get("json")
		value := mapData[0][tag]
		if value == nil && gjangoTag == "required" {
			return errors.New(fmt.Sprintf("field [%s] is required", tag))
		}
	}
	if obj != nil {
		marshal, _ := json.Marshal(mapData)
		_ = json.Unmarshal(marshal, obj)
	}
	return nil
}

func checkStructParams(obj any, decoder *json.Decoder, of reflect.Value) error {
	mapValue := make(map[string]interface{})
	err := decoder.Decode(&mapValue)
	if err != nil {
		return err
	}
	for i := 0; i < of.NumField(); i++ {
		field := of.Type().Field(i)
		name := field.Name
		tag := field.Tag.Get("json")
		// no json tag name
		if tag == "" {
			value := mapValue[name]
			if value == nil {
				return errors.New(fmt.Sprintf("[ERROR] field [%s] does not exist", tag))
			}
		}
		gjangoTag := field.Tag.Get("gjango")
		value := mapValue[tag]
		if value == nil && gjangoTag == "required" {
			return errors.New(fmt.Sprintf("[ERROR] field [%s] does not exist, because [%s] is specified in gjango tag that it is required", tag, tag))
		}
	}
	b, _ := json.Marshal(mapValue)
	_ = json.Unmarshal(b, obj)
	return nil
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
	if code != http.StatusOK {
		c.W.WriteHeader(code)
	}
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
