package server

import (
	"encoding/json"
	"net/http"
    "github.com/benkoben/the-cloud-library/library"
)

type response struct {
	Data    []string `json:"data"`
	Message string   `json:"message"`
}

// newUploadResults creates and returns an uploadRresults.
func newResponse(results []library.Result) response {
	tables := make([]string, 0, len(results))
	for _, result := range results {
		tables = append(tables, string(result.Response))
	}

	return response{
		Data: tables,
		Message: "Success",
	}
}

// JSON returns the JSON encoding of uploadResults.
func (r response) JSON() []byte {
	b, _ := json.Marshal(&r)
	return b
}

// Code returns the status code of uploadResults.
func (r response) Code() int {
	return http.StatusOK
}
