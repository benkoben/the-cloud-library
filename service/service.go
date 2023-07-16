package service

// TODO: Implement a Create operation

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/benkoben/the-cloud-library/library"
)

func main() {
    
    dbCred := library.Credentials{Username: "postgres", Password: "Syp9393"}
    dbOptions := library.PostgresClientOptions{Host: "localhost", SslEnabled: false, Database: "library"}
    serviceOptions := library.ServiceOptions{Timeout: 600, Concurrency: 1}


    dbClient, err := library.NewPgClient(dbCred, dbOptions)
    if err != nil {
        fmt.Errorf("unable to instantiate postgres client")
    }

    library.NewService(dbClient, serviceOptions) 


    r := mux.NewRouter()
    r.HandleFunc("/", HomeHandler)
    r.HandleFunc("/products", ProductsHandler)
    r.HandleFunc("/articles", ArticlesHandler)
    http.Handle("/", r)
}
