package library

// TODO: Each session has its own context and one session is one single connection to the database
// Timeout controls context timout (read about this)

// Concurrency controls how many transactions can be sent to database in parallel

// database implements any driver that implements the InitDb method (this project uses postgres)

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"unicode/utf8"
)


type Result struct {
    table string
    result sql.Result
}

// business logic
type database interface {
    Close() error 
}

type Service interface {
    Create(payloads ...[]byte) ([]Result, error)
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

type operation struct {
    result Result
    operation string
    err error
}

func (s service) Create(table string, payloads ...[]byte)([]Result, error){
    if utf8.RuneCountInString(table) == 0 {
        return nil, errors.New("table cannot be 0 of length")
    } 

    createCh := make(chan []byte)

	// Create a go routine that sends *File to the *File channel.
	// The go routine is needed her to not block the function
	// from continuing to the next steps of calling the
	// producer and consumer methods.
	go func() {
		for _, payload := range payloads {
			createCh <- payload
		}
		// After all payloads has been sent to the channel,
		// close it to signal that no more files will be
		// sent.
		close(createCh)
	}()

    createUploadCh := s.createProducer(table, createCh)  
	return s.createConsumer(createUploadCh)
}


func (s service) createProducer(table string, createCh <-chan []byte) <- chan operation {
    // Add books to the database and save results in a channel
    createUploadCh := make(chan operation)
    var wg sync.WaitGroup

    // table signals which which entity the payloads belong to.
    for i := 1; i <= s.concurreny; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            switch table {
                case "book":
                    for book := range createCh {
                        ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
				        defer cancel()
                       
                        o := operation{}
                        o.operation = "create"

                        b, err := NewBook(book)
                        if err != nil {
                            o.err = err
                        }

                        result, err := b.Add(ctx)
                        if err != nil {
                            o.err = err
                        }
                        o.result = result

                        createUploadCh <- o
                    }
            }
        }()
    }

    // Create an additional go routine that will await all wait groups
	// and close the createUploadCh channel.
	go func() {
		wg.Wait()
		close(createUploadCh)
	}()

    return createUploadCh
}

func (s service) createConsumer(createUploadCh <-chan operation)([]Result, error){
    results := make([]Result, 0, 0)
    errs := make([]error, 0, 0)

    for update := range createUploadCh {
        if update.err != nil {
            errs = append(errs, update.err)
        } else {
            results = append(results, update.result)
        }
    } 
    return results, errors.Join(errs...)
}
