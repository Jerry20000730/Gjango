package Render

import (
	"encoding/json"
	"github.com/Jerry20000730/Gjango/web/Constant"
	"net/http"
)

// JSONRender is a struct that encapsulates data to be rendered as JSON.
type JSONRender struct {
	// Data is the payload to be serialized into JSON.
	Data any
}

// Render writes the JSONRender's Data field to the http.ResponseWriter as JSON.
// It sets the Content-Type header to application/json and serializes the Data field into JSON.
// If the serialization fails, it returns an error.
func (r JSONRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)                  // Set the Content-Type header.
	jsonBytes, err := json.Marshal(r.Data) // Serialize the Data field into JSON.
	if err != nil {
		return err // Return serialization error.
	}
	_, err = w.Write(jsonBytes) // Write the JSON to the response.
	return err                  // Return any error from writing the response.
}

// WriteContentType sets the Content-Type header for the response to JSON_HEADER_CONTENT_TYPE.
func (r JSONRender) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, Constant.JSON_HEADER_CONTENT_TYPE) // Helper function to set Content-Type.
}
