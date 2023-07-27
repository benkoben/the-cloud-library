package server

import (
	"github.com/benkoben/the-cloud-library/library"
	"io/ioutil"
	"net/http"
)

// reponse wraps around the method JSON and Code
type responses interface {
	JSON() []byte
	Code() int
}

// write response to the client.
func write(w http.ResponseWriter, response responses) {
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
				s.log.Println(err)
				write(w, newError(http.StatusBadRequest, errInternalServer))
			}

            // Marshal bytes into books struct
			books, err := library.NewBooks(body)
			if err != nil { 
                s.log.Printf("Handler: bookHandler: NewBooks: %v\n", err)
				write(w, newError(http.StatusBadRequest, errInternalServer))
                return
			} 

            result, err := s.service.StoreBook(*books)
            if err != nil {
                s.log.Printf("Handler: bookHandler: StoreBook: %v\n", err)
				write(w, newError(http.StatusBadRequest, errInternalServer))
                return
            }

            write(w, newResponse(result))
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

func (s *server) healthzHandler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            // read from healthz channel
        } else {
			write(w, newError(http.StatusMethodNotAllowed, "Method not allowed"))
        }
    })
}
