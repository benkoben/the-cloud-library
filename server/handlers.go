package server

import (
	"net/http"
    "io"
    "bytes"
	"github.com/benkoben/the-cloud-library/library"
)

// reponse wraps around the method JSON and Code
type response interface {
    JSON() []byte
    Code() int
}

// write response to the client.
func write(w http.ResponseWriter, response response) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(response.Code())
	w.Write(response.JSON())
}

// Receives one or more books
func (s *server)bookHandler() http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {


        if r.Method == http.MethodPost {
            // TODO: ioutil.Reader() can be used here to read the body (https://www.digitalocean.com/community/tutorials/how-to-make-an-http-server-in-go#reading-a-request-body)
            book, err := library.NewBook(streamToByte(r.Body))
            if err != nil {
                write(w, newError(http.StatusBadRequest, errInternalServer))
			    return
            }
            s.service.book.Store(r.Context(), book)
        }
        if r.Method == http.MethodGet {

        }
        if r.Method == http.MethodPut {

        }
        if r.Method == http.MethodDelete {

        }
    })
}


