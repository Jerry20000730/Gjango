package Render

import "net/http"

// Render is an interface that defines methods for rendering content and setting
// the content type of HTTP responses.
type Render interface {
	// Render writes the rendered content to the http.ResponseWriter and returns
	// an error if the rendering fails.
	Render(w http.ResponseWriter) error

	// WriteContentType sets the Content-Type header of the HTTP response.
	WriteContentType(w http.ResponseWriter)
}

// writeContentType is a utility function that sets the Content-Type header of
// the HTTP response to the specified value.
func writeContentType(w http.ResponseWriter, value string) {
	w.Header().Set("Content-Type", value)
}
