package server

import (
    "encoding/json"
)

const (
	errInternalServer       = "Internal server error."
	errMethodNotAllowed     = "Method not allowed."
//	errNotMultipartError    = "Malformed request. Not multipart/form-data."
    errMissingFieldBook     = "Malformed request. Request body cannot be marshaled into Book"
	errUnauthorized         = "Not authorized."
)

// Error represents an HTTP error response from the server.
type Error struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

// newError creates and returns an Error.
func newError(statusCode int, message string) Error {
	return Error{
		StatusCode: statusCode,
		Message:    message,
	}
}

// Error implements interface error.
func (e Error) Error() string {
	return e.Message
}

// JSON returns the JSON encoding of Error.
func (e Error) JSON() []byte {
	b, _ := json.Marshal(&e)
	return b
}

// Code returns the status code of the Error.
func (e Error) Code() int {
	return e.StatusCode
}
