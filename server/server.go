package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/benkoben/the-cloud-library/library"
)

const (
    defaultPort = 3000
    defaultReadTimeout  = time.Second * 15
    defaultWriteTimeout = time.Second * 15
    defaultIdleTimeout  = time.Second * 15
    defaultMultiPartMaxBytes = 32 << 20
)

type logger interface {
	Printf(format string, v ...any)
	Println(v ...any)
	Fatalf(format string, v ...any)
	Fatalln(v ...any)
}

// server is the serving part of the application containing all handlers
// and services.
type server struct {
	httpServer        *http.Server
	router            *http.ServeMux
	log               logger
	service           library.Service
}

// Options contains options for the server.
type Options struct {
	Router            *http.ServeMux
	Service           library.Service
	Log               logger
	Host              string
	Port              int
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

func NewServer(options Options) (*server, error){
    
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
   
    srv := &http.Server{
        Addr: options.Host + ":" + strconv.Itoa(options.Port),
        Handler: options.Router,
        ReadTimeout: options.ReadTimeout,
        WriteTimeout: options.WriteTimeout,
        IdleTimeout: options.IdleTimeout,
    }
    return &server{
       httpServer: srv,
       router: options.Router,
       log: options.Log,
       service: options.Service,
    }, nil
}

func (s *server) Start() {
    s.routes()
    go func() {
         
    }()
}
