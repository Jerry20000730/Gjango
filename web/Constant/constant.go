package Constant

import "net/http"

const ANY = "ANY"
const GET = http.MethodGet
const POST = http.MethodPost
const PUT = http.MethodPut
const DELETE = http.MethodDelete
const PATCH = http.MethodPatch
const HEAD = http.MethodHead
const OPTIONS = http.MethodOptions

// HTML_HEADER_CONTENT_TYPE defines the Content-Type header for HTML responses.
const HTML_HEADER_CONTENT_TYPE = "text/html; charset=utf-8"

// JSON_HEADER_CONTENT_TYPE is a constant for the Content-Type header for JSON responses.
const JSON_HEADER_CONTENT_TYPE = "application/json; charset=utf-8"

// STRING_HEADER_CONTENT_TYPE defines the Content-Type header for plain text responses.
const STRING_HEADER_CONTENT_TYPE = "text/plain; charset=utf-8"

// XML_HEADER defines the Content-Type header for XML responses.
const XML_HEADER = "application/xml; charset=utf-8"

const DEFAULT_MAX_MEMORY = 32 << 20 // 32 MB
