package library

import (
	"context"
	"errors"
	"sync"
	"time"
)

type dbClient interface {
    IsHealthy(context.Context) bool
}

type Result struct {
	Response []byte
}

// tablesStores wraps all CRUD operations for a table inside the database
type bookStore interface {
	Store(context.Context, *Book) (Result, error)
	Get(context.Context, int64) (Result, error)
	Delete(context.Context, *Book) (Result, error)
	List(context.Context, *BooksFilters) (Result, error)
}

// Each table in the datbase has its own tableStore.
type DbStore struct {
	Books   bookStore
	// Users   UserStore
	// Rentals RentalStore
}

// service is used to control the behaviour our service should
// have towards the database
type Service struct {
	Store      DbStore
    Client     dbClient
	Timeout    time.Duration
	Concurreny int
}

type ServiceOptions struct {
	Timeout     time.Duration
	Concurrency int
}

// Constructor function for the service
func NewService(client dbClient, store DbStore, options ServiceOptions) (*Service, error) {
	if store.Books == nil {
		return nil, errors.New("store.books must not be nil")
	}

//	if store.Users == nil {
//		return nil, errors.New("store.users must not be nil")
//	}
//
//	if store.Rentals == nil {
//		return nil, errors.New("store.rentals must not be nil")
//	}

    if client == nil {
        return nil, errors.New("client cannot be nil")
    }

	if options.Concurrency == 0 {
		options.Concurrency = 1
	}

	if options.Timeout == 0 {
		options.Timeout = time.Second * 10
	}

	return &Service{
		Store:      store,
        Client:     client,
		Timeout:    options.Timeout,
		Concurreny: options.Concurrency,
	}, nil
}

type operation struct {
	result    Result
	operation string
	err       error
}

func (s Service) StoreBook(books Books) ([]Result, error) {

	storeBookCh := make(chan *Book)

	// Create a go routine that sends *payload to the *storeCh channel.
	// The go routine is needed here to not block the function
	// from continuing to the next steps of calling the
	// producer and consumer methods.
	go func() {
		for _, book := range books.Data {
			storeBookCh <- &book
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

				res, err := s.store.Books.Store(ctx, book)
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
