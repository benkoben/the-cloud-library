package main

import (
    "log"
    "net/http"

    "github.com/benkoben/the-cloud-library/server"
    "github.com/benkoben/the-cloud-library/config"

    "fmt"
)

func main() {
	// get config
    cfg, err := config.New()
    if err != nil {
        fmt.Errorf("could not create config: %s\n", err)
    }
    
    // new service
    svc, err := config.NewLibraryService(&cfg.Library)
    if err != nil {
        fmt.Errorf("could not create library service: %s\n", err)
    }


    srv := server.New(server.Options{
		Router:            http.NewServeMux(),
        Service:           svc,
		Log:               log.Default(),
		Host:              cfg.Server.Host,
		Port:              cfg.Server.Port,
		ReadTimeout:       cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		IdleTimeout:       cfg.Server.IdleTimeout,
	})
}

