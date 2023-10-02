package main

import (
    "log"

    "github.com/benkoben/the-cloud-library/server"
    "github.com/benkoben/the-cloud-library/config"

    "github.com/gorilla/mux"

)

func main() {
	// get config
    cfg, err := config.New()
    if err != nil {
        log.Fatalf("could not create config: %s\n", err)
    }
    
    // new service
    svc, err := config.NewLibraryService(&cfg.Library)
    if err != nil {
        log.Fatalf("could not create library service: %s\n", err)
    }


    srv, err := server.New(server.Options{
		Router:            mux.NewRouter(),
        Service:           *svc,
		Log:               log.Default(),
		Host:              cfg.Server.Host,
		Port:              cfg.Server.Port,
		ReadTimeout:       cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		IdleTimeout:       cfg.Server.IdleTimeout,
	})

    if err != nil {
        log.Fatalf("could not start server: %s\n", err)
    }

    srv.Start()

}

