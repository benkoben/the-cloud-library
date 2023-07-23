package library

// TODO: Each session has its own context and one session is one single connection to the database
// Timeout controls context timout (read about this)

// Concurrency controls how many transactions can be sent to database in parallel

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Result struct {
	Response []byte
}

type tableStore interface {
	Store(context.Context, any) (Result, error)
	Get(context.Context, int64) error
	Delete(context.Context, any) error
	List(context.Context, any) error
}

type DbStore struct {
	books   tableStore
	users   tableStore
	rentals tableStore
}

// service is used to control the behaviour our service should
// have towards the database
type Service struct {
	store      DbStore
	timeout    time.Duration
	concurreny int
}

type ServiceOptions struct {
	Timeout     time.Duration
	Concurrency int
}

// Constructor function for the service
func NewService(store DbStore, options ServiceOptions) (*Service, error) {
	if store.books == nil {
		return nil, errors.New("store.books must not be nil")
	}

	if store.users == nil {
		return nil, errors.New("store.users must not be nil")
	}

	if store.rentals == nil {
		return nil, errors.New("store.rentals must not be nil")
	}

	if options.Concurrency == 0 {
		options.Concurrency = 1
	}

	if options.Timeout == 0 {
		options.Timeout = time.Second * 10
	}

	return &Service{
		store:      store,
		timeout:    options.Timeout,
		concurreny: options.Concurrency,
	}, nil
}

type operation struct {
	result    Result
	operation string
	err       error
}

func (s Service) StoreBook(payloads ...[]byte) ([]Result, error) {
	if len(payloads) == 0 {
		return nil, fmt.Errorf("payload cannot be 0 in length")
	}

	storeBookCh := make(chan *Book)

	// Create a go routine that sends *payload to the *storeCh channel.
	// The go routine is needed here to not block the function
	// from continuing to the next steps of calling the
	// producer and consumer methods.
	go func() {
		for _, payload := range payloads {
			if b, err := NewBook(payload); err != nil {
				storeBookCh <- b
			}
		}
		// After all payloads has been sent to the channel,
		// close it to signal that no more files will be
		// sent.
		close(storeBookCh)
	}()

	storeBookResultCh := s.storeBookProducer(storeBookCh)
	return s.storeBookResultConsumer(storeBookResultCh)
}

func (s Service) storeBookProducer(storeBookCh <-chan *Book) <-chan operation {
	// Add books to the database and save results in a channel
	storeBookResultCh := make(chan operation)
	var wg sync.WaitGroup

	// table signals which which entity the payloads belong to.
	for i := 1; i <= s.concurreny; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for book := range storeBookCh {
				ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
				defer cancel()

				o := operation{}
				o.operation = "store"

				res, err := s.store.books.Store(ctx, &book)
				if err != nil {
					o.err = err
				}
				o.result = res

				storeBookResultCh <- o
			}
		}()
	}

	// Create an additional go routine that will await all wait groups
	// and close the createUploadCh channel.
	go func() {
		wg.Wait()
		close(storeBookResultCh)
	}()

	return storeBookResultCh
}

func (s Service) storeBookResultConsumer(storeBookResultCh <-chan operation) ([]Result, error) {
	results := make([]Result, 0, 0)
	errs := make([]error, 0, 0)

	for res := range storeBookResultCh {
		if res.err != nil {
			errs = append(errs, res.err)
		} else {
			results = append(results, res.result)
		}
	}
	return results, errors.Join(errs...)
}
