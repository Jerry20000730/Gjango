package Render

import (
	"encoding/xml"
	"github.com/Jerry20000730/Gjango/web/Constant"
	"net/http"
)

// XMLRender is a renderer for XML data.
type XMLRender struct {
	// Data holds the actual data to be rendered as XML.
	Data any
}

// Render writes the XML representation of Data to the http.ResponseWriter.
// It sets the Content-Type header using WriteContentType before encoding Data as XML.
// Returns an error if encoding fails.
func (r XMLRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	return xml.NewEncoder(w).Encode(r.Data)
}

// WriteContentType sets the Content-Type header of the http.ResponseWriter to XML_HEADER.
func (r XMLRender) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, Constant.XML_HEADER)
}
