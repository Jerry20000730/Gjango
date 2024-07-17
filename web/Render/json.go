package Render

import (
	"encoding/json"
	"net/http"
)

type JSONRender struct {
	Data any
}

const JSON_HEADER_CONTENT_TYPE = "application/json; charset=utf-8"

func (r JSONRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	jsonBytes, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}

func (r JSONRender) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, JSON_HEADER_CONTENT_TYPE)
}
