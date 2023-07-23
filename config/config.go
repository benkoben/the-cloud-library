package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/benkoben/the-cloud-library/library"
	"github.com/caarlos0/env/v6"
)

// Default settings for the server
const (
	defaultHost         = "0.0.0.0"
	defaultPort         = 3000
	defaultReadTimout   = time.Second * 15
	defaultWriteTimeout = time.Second * 15
	defaultIdleTimeout  = time.Second * 30
)

// Default settings for the library service
const (
	defaultDatabaseHost        = "127.0.0.1"
	defaultDatabasePort        = "5432"
	defaultDatabaseSslEnabled  = false
	defaultDatabaseName        = "library"
	defaultDatabaseCredentials = "/credentials.json"
	defaultLibraryTimeout      = time.Second * 30
	defaultLibraryConcurrency  = 5
)

// Configuration defines all settings for the whole application
type Configuration struct {
	Server  Server
	Library Library
}

// Server defines all the settings for the server component of the application
type Server struct {
	Host         string        `env:"LIBRARY_LISTEN_HOST"`
	Port         int           `env:"LIBRARY_LISTEN_PORT"`
	ReadTimeout  time.Duration `env:"LIBRARY_READ_TIMEOUT"`
	WriteTimeout time.Duration `env:"LIBRARY_WRITE_TIMEOUT"`
	IdleTimeout  time.Duration `env:"LIBRARY_IDLE_TIMEOUT"`
}

// Librabry defines all the settings for the database service component of the application
type Library struct {
	DatabaseHost        string        `env:"LIBRARY_SERVICE_DB_HOST"`
	DatabaseSslEnabled  bool          `env:"LIBRARY_SERVICE_DB_SSL_ENABLED"`
	DatabasePort        string        `env:"LIBRARY_SERVICE_DB_PORT"`
	DatabaseName        string        `env:"LIBRARY_SERVICE_DB_NAME"`
	DatabaseCredentials string        `env:"LIBRARY_SERVICE_DB_CRED_PATH"`
	Timeout             time.Duration `env:"LIBRARY_SERVICE_TIMEOUT"`
	Concurrency         int           `env:"LIBRARY_CONCURRENCY"`
}

// Creates a new configuration for the application. Which can be used to start the server
// and connect to the database layer.
//
// Defaults have lowest priority and are overwritten by environment variables
func New() (*Configuration, error) {
	cfg := &Configuration{
		Server: Server{
			Host:         defaultHost,
			Port:         defaultPort,
			ReadTimeout:  defaultReadTimout,
			WriteTimeout: defaultWriteTimeout,
			IdleTimeout:  defaultIdleTimeout,
		},
		Library: Library{
			DatabaseHost:        defaultDatabaseHost,
			DatabasePort:        defaultDatabasePort,
			DatabaseSslEnabled:  defaultDatabaseSslEnabled,
			DatabaseName:        defaultDatabaseName,
			DatabaseCredentials: defaultDatabaseCredentials,
			Timeout:             defaultLibraryTimeout,
			Concurrency:         defaultLibraryConcurrency,
		},
	}
    // overwrite defaults with environment variables.
    // if no environment variable is found then the default value is kept.
	if err := env.Parse(cfg); err != nil {
		return nil, errors.New("configuration: could not pase environment variables")
	}
	return cfg, nil
}

// Tie all settings for the database layer together and returns a Service
func NewLibraryService(cfg *Library) (*Library.Service, error) {
	cred, err := library.NewDbCredentials(cfg.DatabaseCredentials)
	if err != nil {
		return nil, fmt.Errorf("cannot create db credentials: %s\n", err)
	}

	dbOptions := &library.PostgresClientOptions{
		Host:       cfg.DatabaseHost + ":" + cfg.DatabasePort,
		SslEnabled: cfg.DatabaseSslEnabled,
		Database:   cfg.DatabaseName,
	}

	dbClient, err := library.NewPgClient(cred, *dbOptions)
	if err != nil {
		return nil, fmt.Errorf("cannot create new db client: %s\n", err)
	}

    // bookStore implements all CRUD operations for the books table
	bookStore, err := library.NewBookStore(dbClient)
	if err != nil {
		return nil, fmt.Errorf("could not create bookStore: %s\n", err)
	}

    // userStore implements all CRUD operations for the books table
	userStore, err := library.NewUserStore(dbClient)
	if err != nil {
		return nil, fmt.Errorf("could not create userStore: %s\n", err)
	}

    // rentalStore implements all CRUD operations for the books table
	rentalStore, err := library.NewRentalStore(dbClient)
	if err != nil {
		return nil, fmt.Errorf("could not create rentalStore: %s\n", err)
	}

	dbStore = library.DbStore{
		books:       bookStore,
		users:       userStore,
		rentalStore: rentalStore,
	}

	opts := library.ServiceOptions{
		Timeout:     cfg.Timeout,
		Concurrency: cfg.Concurrency,
	}

	return library.NewService(dbStore, opts)
}
