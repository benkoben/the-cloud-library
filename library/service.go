package library

import (
	"errors"
	"time"
    "database/sql"
)

type Env struct {
    db *sql.DB
}

type Result struct {
    Message string
}

// business logic
type database interface {
    InitDb() error
}

// service is used to control the behaviour our service should
// have towards the database
type service struct {
	database   database
	timeout    time.Duration
	concurreny int
}

type ServiceOptions struct {
	Timeout     time.Duration
	Concurrency int
}

// Constructor function for the service
func NewService(database database, options ServiceOptions) (*service, error) {
	if database == nil {
		return nil, errors.New("database must not be nil")
	}

	if options.Concurrency == 0 {
		options.Concurrency = 1
	}

	if options.Timeout == 0 {
		options.Timeout = time.Second * 10
	}

	return &service{
		database:   database,
		timeout:    options.Timeout,
		concurreny: options.Concurrency,
	}, nil
}


// TODO stuff below.
// CRUD methods should be able to detect which table they should operate on and prepare channels accordingly
func (s service) Create(table string)([]Result, error){return nil, nil}
func (s service) Remove(table string)([]Result, error){return nil, nil}
func (s service) Update(table string)([]Result, error){return nil, nil}
func (s service) Delete(table string)([]Result, error){return nil, nil}


func (s service) createProducer(table string)([]Result, error){return nil, nil}
func (s service) removeProducer(table string)([]Result, error){return nil, nil}
func (s service) updateProducer(table string)([]Result, error){return nil, nil}
func (s service) deleteProducer(table string)([]Result, error){return nil, nil}

func (s service) createConsumer(table string)([]Result, error){return nil, nil}
func (s service) removeConsumer(table string)([]Result, error){return nil, nil}
func (s service) updateConsumer(table string)([]Result, error){return nil, nil}
func (s service) deleteConsumer(table string)([]Result, error){return nil, nil}
