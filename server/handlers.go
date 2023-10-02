package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/benkoben/the-cloud-library/library"

	"github.com/gorilla/mux"
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
        // Use http.MaxBytesReader to enforce a maximum read of 1MB from the
        // response body. A request body larger than that will now result in
        // Decode() returning a "http: request body too large" error.
        r.Body = http.MaxBytesReader(w, r.Body, 1048576)

		if r.Method == http.MethodPost {

            var books []library.Book


            body, err := io.ReadAll(r.Body)
            if err != nil {
                s.log.Fatalf("Handler: could not read body: %v", err)
            }
            // Setup the decoder and call the DisallowUnknownFields() method on it.
            // This will cause Decode() to return a "json: unknown field ..." error
            // if it encounters any extra unexpected fields in the JSON. Strictly
            // speaking, it returns an error for "keys which do not match any
            // non-ignored, exported fields in the destination".
            dec := json.NewDecoder(bytes.NewReader(body))
            dec.DisallowUnknownFields()

        	for {
        		var book library.Book

        		if err := dec.Decode(&book); err == io.EOF {
        			break
        		} else if err != nil {
                    s.log.Fatalf("Handler: failed to read from buffer: %v\n", err)
        		}
                
                books = append(books, book)    
        	}


            result, err := s.service.StoreBook(books)
            if err != nil {
                s.log.Printf("Handler: bookHandler: StoreBook: %v\n", err)
				write(w, newError(http.StatusBadRequest, errInternalServer))
                return
            }

            write(w, newResponse(result))
		}
		// TODO: implement the following methods
		if r.Method == http.MethodGet {
            vars := mux.Vars(r)
            id, ok := vars["id"]
            if !ok {
                s.log.Println("Handler: bookHandler: missing id parameter in request url")
                write(w, newError(http.StatusBadRequest, errMissingParameter))
            }
            // The request should be routed to a service component
            // that retrieves the requested book. This component does not have to
            // be concurrent because we are only requesting a single book.
            id64, err := strconv.Atoi(id)
            if err != nil {
                s.log.Printf("Handler: bookHandler: %v\n", err)
                write(w, newError(http.StatusInternalServerError, errInternalServer))
            }
            result, err := s.service.GetBook(int64(id64))

            write(w, newResponse(result))
		}
		if r.Method == http.MethodPut {
			write(w, newError(http.StatusMethodNotAllowed, "Method not allowed"))
		}
		if r.Method == http.MethodDelete {
			write(w, newError(http.StatusMethodNotAllowed, "Method not allowed"))
		}
	})
}
