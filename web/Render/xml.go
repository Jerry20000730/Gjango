package Render

import (
	"encoding/xml"
	"net/http"
)

const XML_HEADER = "application/xml; charset=utf-8"

type XMLRender struct {
	Data any
}

func (r XMLRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	return xml.NewEncoder(w).Encode(r.Data)
}

func (r XMLRender) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, XML_HEADER)
}
