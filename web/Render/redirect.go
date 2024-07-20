package Render

import (
	"fmt"
	"net/http"
)

// Redirect represents an HTTP redirect response.
// It holds the status code, the original request, and the location to redirect to.
type Redirect struct {
	// HTTP status code for the redirect
	Code int
	// The original request
	Request *http.Request
	// The URL to redirect to
	Location string
}

// Render executes the redirect by writing the HTTP status code and the "Location" header to the response writer.
// It panics if the status code is not within the range of valid redirect codes (300-308) except for 201 (Created).
// This method satisfies the Renderer interface, allowing it to be used wherever Renderer is accepted.
func (r Redirect) Render(w http.ResponseWriter) error {
	if (r.Code < http.StatusMultipleChoices || r.Code > http.StatusPermanentRedirect) && r.Code != http.StatusCreated {
		panic(fmt.Sprintf("Cannot redirect with status code %d", r.Code))
	}
	http.Redirect(w, r.Request, r.Location, r.Code)
	return nil
}

// WriteContentType is a no-op for Redirect since redirects do not require a content type.
// This method is required to satisfy the Renderer interface.
func (r Redirect) WriteContentType(http.ResponseWriter) {}
