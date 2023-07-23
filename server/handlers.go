package server

import (
	"github.com/benkoben/the-cloud-library/library"
	"io/ioutil"
	"net/http"
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
func (s *server) bookHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodPost {

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				s.log.Fatalln(err)
				write(w, newError(http.StatusBadRequest, errInternalServer))
			}

			book, err := library.NewBook(body)
			if err != nil {
				s.log.Fatalln(err)
				write(w, newError(http.StatusBadRequest, errMissingFieldBook))
			}

			s.service.books.Store(r.Context(), book)
		}
		// TODO: implement the following methods
		if r.Method == http.MethodGet {
			write(w, newError(http.StatusMethodNotAllowed, "Method not allowed"))
		}
		if r.Method == http.MethodPut {
			write(w, newError(http.StatusMethodNotAllowed, "Method not allowed"))
		}
		if r.Method == http.MethodDelete {
			write(w, newError(http.StatusMethodNotAllowed, "Method not allowed"))
		}
	})
}
