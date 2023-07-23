package server


import (
	"encoding/json"
	"net/http"

	"github.com/benkoben/the-cloud-library/library"
)


type Results struct {
    Stores  []string `json:"stores"`
    Message string `json:"message"`
}

// JSON returns the JSON encoding of uploadResults.
func (r Results) JSON() []byte {
	b, _ := json.Marshal(&r)
	return b
}

// Code returns the status code of uploadResults.
func (r Results) Code() int {
	return http.StatusOK
}
