package Render

import (
	"fmt"
	"github.com/Jerry20000730/Gjango/web/Constant"
	"github.com/Jerry20000730/Gjango/web/Render/internal/bytesconv" // Internal package for byte conversion utilities.
	"net/http"
)

// StringRender is a renderer for string responses. It supports formatting according to the provided format and data.
type StringRender struct {
	// Format is the string format, similar to those used in fmt.Sprintf
	Format string
	// Data contains the values to be used when formatting the string.
	Data []any
}

// Render writes the formatted string or plain string to the response writer.
// It sets the Content-Type header to "text/plain; charset=utf-8".
// If Data is not empty, it formats the string according to Format and Data.
// Otherwise, it writes the Format string as is.
func (r StringRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	if len(r.Data) > 0 {
		_, err := fmt.Fprintf(w, r.Format, r.Data...) // Format string with data and write to response.
		return err
	}
	_, err := w.Write(bytesconv.StringToBytes(r.Format)) // Convert Format string to bytes and write to response.
	return err
}

// WriteContentType sets the Content-Type header for the response.
func (r StringRender) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, Constant.STRING_HEADER_CONTENT_TYPE)
}
