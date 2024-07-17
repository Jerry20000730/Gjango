package Render

import (
	"html/template"
	"net/http"
)

const HTML_HEADER_CONTENT_TYPE = "text/html; charset=utf-8"

type HTMLRender struct {
	Template   *template.Template
	Name       string
	Data       any
	IsTemplate bool
}

type HTMLPreloader struct {
	FuncMap  template.FuncMap
	Template *template.Template
}

func (r HTMLRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	if !r.IsTemplate {
		_, err := w.Write([]byte(r.Data.(string)))
		return err
	}
	err := r.Template.ExecuteTemplate(w, r.Name, r.Data)
	return err
}

func (r HTMLRender) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, HTML_HEADER_CONTENT_TYPE)
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
