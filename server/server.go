package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/benkoben/the-cloud-library/library"
)

const (
	defaultPort              = 3000
	defaultReadTimeout       = time.Second * 15
	defaultWriteTimeout      = time.Second * 15
	defaultIdleTimeout       = time.Second * 15
)

var ( 
)

type status struct {
    Healthy bool `json:'healthy'`
}

type logger interface {
	Printf(format string, v ...any)
	Println(v ...any)
	Fatalf(format string, v ...any)
	Fatalln(v ...any)
}

// server is the serving part of the application containing all handlers
// and services.
type server struct {
	httpServer *http.Server
	router     *http.ServeMux
	log        logger
	service    library.Service
    status status
}

// Options contains options for the server.
type Options struct {
	Router       *http.ServeMux
	Service      library.Service
	Log          logger
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func New(options Options) (*server, error) {

	if options.Router == nil {
		options.Router = http.NewServeMux()
	}
	if options.Port == 0 {
		options.Port = defaultPort
	}
	if options.Host == "" {
		options.Host = "0.0.0.0"
	}
	if options.ReadTimeout == 0 {
		options.ReadTimeout = defaultReadTimeout
	}
	if options.WriteTimeout == 0 {
		options.WriteTimeout = defaultWriteTimeout
	}
	if options.IdleTimeout == 0 {
		options.IdleTimeout = defaultIdleTimeout
	}

    status := status{Healthy: true}

	srv := &http.Server{
		Addr:         options.Host + ":" + strconv.Itoa(options.Port),
		Handler:      options.Router,
		ReadTimeout:  options.ReadTimeout,
		WriteTimeout: options.WriteTimeout,
		IdleTimeout:  options.IdleTimeout,
	}

	return &server{
		httpServer: srv,
		router:     options.Router,
		log:        options.Log,
		service:    options.Service,
        status: status,
	}, nil
}


// Start the server.
func (s server) Start() {
	s.routes()
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Fatalf("Server error %v.\n", err)
		}
		s.log.Println("Server stopped.")
	}()

    // implement go routine that checks health
    go func() {
            time.Sleep("")
    }()

	s.log.Printf("Server listening on: %s.\n", s.httpServer.Addr)
	s.stop()
}

func (s server) stop() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop

	s.log.Printf("Shutting down server. Reason: %s.\n", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	s.httpServer.SetKeepAlivesEnabled(false)
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.log.Printf("Server shutdown: %v.\n", err)
	}
}
