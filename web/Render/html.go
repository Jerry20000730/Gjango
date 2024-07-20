// Package Render provides functionality to render HTML content.
package Render

import (
	"github.com/Jerry20000730/Gjango/web/Constant"
	"html/template"
	"net/http"
)

// HTMLRender is a struct for rendering HTML templates or strings.
type HTMLRender struct {
	Template   *template.Template // Template is the HTML template to be rendered.
	Name       string             // Name is the template name for lookup in case of multiple templates.
	Data       any                // Data contains the data to be passed to the template.
	IsTemplate bool               // IsTemplate indicates whether to render a template or a plain string.
}

// HTMLPreloader is a struct to preload templates with specific functions.
type HTMLPreloader struct {
	FuncMap  template.FuncMap   // FuncMap is a map of functions to be used in templates.
	Template *template.Template // Template is the preloaded template.
}

// Render writes the rendered template or string to the http.ResponseWriter.
// If IsTemplate is false, it writes Data as a plain string.
// Otherwise, it executes the template with the provided data.
func (r HTMLRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w) // Set the Content-Type header.
	if !r.IsTemplate {
		_, err := w.Write([]byte(r.Data.(string))) // Write the plain string.
		return err
	}
	err := r.Template.ExecuteTemplate(w, r.Name, r.Data) // Execute the template.
	return err
}

// WriteContentType sets the Content-Type header for the response.
func (r HTMLRender) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, Constant.HTML_HEADER_CONTENT_TYPE)
}

//func (r *HTMLRender) PreLoadTemplate(ctx *context.Context, statusCode int, name string, data any) error {
//	if r.GetTemplate() == nil {
//		return fmt.Errorf("no pre-load template has been found, have you load the template in the engine")
//	}
//	ctx.W.WriteHeader(statusCode)
//	ctx.W.Header().Set("Content-Type", "text/html; charset=utf-8")
//	err := r.GetTemplate().ExecuteTemplate(ctx.W, name, data)
//	if err != nil {
//		return err
//	}
//	return err
//}
//
//func (r *HTMLRender) Template(ctx *context.Context, statusCode int, name string, data any, filenames ...string) error {
//	ctx.W.WriteHeader(statusCode)
//	ctx.W.Header().Set("Content-Type", "text/html; charset=utf-8")
//	t := template.New(name)
//	t, err := t.ParseFiles(filenames...)
//	if err != nil {
//		return err
//	}
//	err = t.Execute(ctx.W, data)
//	return err
//}
//
//func (r *HTMLRender) TemplateGlob(ctx *context.Context, statusCode int, name string, data any, pattern string) error {
//	ctx.W.WriteHeader(statusCode)
//	ctx.W.Header().Set("Content-Type", "text/html; charset=utf-8")
//	t := template.New(name)
//	t, err := t.ParseGlob(pattern)
//	if err != nil {
//		return err
//	}
//	err = t.Execute(ctx.W, data)
//	return err
//}
