package server

import (
	"encoding/json"
	"net/http"
)


type response struct {
    Res interface{} `json:"res"`
} 

// newUploadResults creates and returns an uploadRresults.
func newResponse(results any) response {
	return response{
		Res: results,
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
