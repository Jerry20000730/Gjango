package Render

import (
	"fmt"
	"github.com/Jerry20000730/Gjango/web/Render/internal/bytesconv"
	"net/http"
)

const STRING_HEADER_CONTENT_TYPE = "text/plain; charset=utf-8"

type StringRender struct {
	Format string
	Data   []any
}

func (r StringRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	if len(r.Data) > 0 {
		_, err := fmt.Fprintf(w, r.Format, r.Data...)
		return err
	}
	_, err := w.Write(bytesconv.StringToBytes(r.Format))
	return err
}

func (r StringRender) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, STRING_HEADER_CONTENT_TYPE)
}
