package library

// Here we will be using a SDK to communicatate with postgres
// Lets make a struct that has all the CRUD methods that full
// the dbOperator struct

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
)

var (
    ErrNotFound = errors.New("not found")
)

type Books struct {
    Response []Book `json:"response"`
}

type Book struct {
    Id             int        `json:"id" validate:"required"`
	Isbn           string     `json:"isbn" validate:"required"`
	Title          string     `json:"title" validate:"required"`
	Lang           string     `json:"lang" validate:"required"`
	Translator     string     `json:"translator"`
	Author         string     `json:"author" validate:"required"`
	Pages          int        `json:"pages" validate:"required"`
	Publisher      string     `json:"publisher" validate:"required"`
	Published_date *time.Time `json:"published_date" validate:"required"`
	Added_date     *time.Time `json:"added_date" validate:"required"`
}

func (b *Books) Marshal() ([]byte, error) {
    return json.Marshal(b)
}

func NewBook(payload []byte) (*Book, error) {
    if len(payload) == 0 {
        return nil, fmt.Errorf("payload cannot be 0 in size")
    }
    var book Book
    if err := json.Unmarshal(payload, &book); err != nil {
        return nil, fmt.Errorf("cannot unmarshall payload into book %w,", err)
    }
    return &book, nil
}

type BookStore struct {
    db *sql.DB
}

// Constructor method used to instantiate a new BookStore
func NewBookStore(db *sql.DB) (*BookStore, error) {
    if db == nil {
        return nil, errors.New("db cannot be nil")
    }
    return &BookStore{db: db}, nil
}

// Store saves a book to the database. If the book has no ID then it will be updated. Otherwise,
// it will be inserted and the ID will be set.
//
// If the book has an ID and it does not exist in the database, Store returns ErrNotFound.
func (bs *BookStore) Store(ctx context.Context, b *Book) (Result, error) {
    r := Books{Response: []Book{*b}}

    response, err := r.Marshal()
    if err != nil {
        return Result{}, fmt.Errorf("could not marshal book %w,", err)
   }
    res := Result{Response: response}

    if b.Id == 0 {
        return res, bs.insert(ctx, b)
    }

    return res, bs.update(ctx, b)
}

// Retrieves a specific book from the database based on the id argument
func (bs *BookStore) Get(ctx context.Context, id int64) (Result, error) {
    var b Book 
    err := squirrel.
            Select("id", "isbn", "title", "translator", "author", "pages", "publisher", "published_date", "added_date").
            From("books").
            Where("id = ?", id).
            RunWith(bs.db).
            QueryRowContext(ctx).
            Scan(&b.Id, &b.Isbn, &b.Title, &b.Translator, &b.Author, &b.Pages, &b.Publisher, &b.Published_date, &b.Added_date)

    if err != nil {
        return Result{}, err
    }
   
    // Build response message
    r := Books{Response: []Book{b}}

    response, err := r.Marshal()
    if err != nil {
        return Result{}, fmt.Errorf("could not marshal book %w,", err)
    }


    return Result{
        Response: response,
    }, nil
}

// Add a book to the books table
func (bs *BookStore) insert(ctx context.Context, b *Book) error {
    return squirrel.
            Insert("books").
            Columns("id", "isbn", "title", "translator", "author", "pages", "publisher", "published_date", "added_date").
            Values(b.Id, b.Isbn, b.Title, b.Translator, b.Author, b.Pages, b.Publisher, b.Published_date, b.Added_date).
            Suffix("RETURNING id").
            RunWith(bs.db).
            QueryRowContext(ctx).
            Scan(&b.Id)
}

// Altters the rows for a specific book
//
// If no rows where updated then a ErrNotFound is returned
func (bs *BookStore) update(ctx context.Context, b *Book) error {
    res, err := squirrel.
            Update("books").
            Set("id", b.Id).
            Set("isbn", b.Isbn).
            Set("isbn", b.Isbn).
            Set("title", b.Title).
            Set("translator", b.Translator).
            Set("author", b.Author).
            Set("pages", b.Pages).
            Set("publisher", b.Publisher).
            Set("published_date", b.Published_date).
            Set("added_date", b.Added_date).
            Where("id = ?", b.Id).
            RunWith(bs.db).
            ExecContext(ctx)

    if err != nil {
        return fmt.Errorf("update books: %w", err)
    }
    rows, _ := res.RowsAffected()
    if rows == 0 {
        return ErrNotFound 
    }
    return nil
}

// Delete removes a book from the database
//
// If the delete operation returns an empty reponse then a ErrNotFound is returned
func (bs *BookStore) Delete(ctx context.Context, b *Book) (Result, error) {

    res, err := squirrel.
                Delete("books").
                Where("id = ? ", b.Id).
                RunWith(bs.db).
                ExecContext(ctx)

    if err != nil {
        return Result{}, fmt.Errorf("could not marshal book %w,", err)
    }

    rows, _ := res.RowsAffected()


    if rows == 0 {
        return Result{}, ErrNotFound
    }

    r := Books{Response: []Book{*b}}
    
    // Create response message
    response, err := r.Marshal()
    if err != nil {
        return Result{}, fmt.Errorf("could not marshal book %w,", err)
    }

    return Result{
        Response: response,
    }, nil
}



type BooksFilters struct {
    // Id matches a books ID
    Id             int
    // Isbn matches a books isbn number
	Isbn           string
    // Title matches a books Title
	Title          string
    // Lang matches all books published in a certain language
	Lang           string
    // Translator matches all books translated by a certain Translator
	Translator     string
    // Author matches all books written by a certain Author
	Author         string
    // Publisher matches all books written by a certain Publisher
	Publisher      string
}

// List searches for characters in the database.
//
// If filters is nil, all characters are returned. Otherwise, the results are
// filtered by the criteria in filters.
func (bs *BookStore) List(ctx context.Context, filters *BooksFilters) (Result, error) {
    	q := squirrel.
		    Select("b.id", "b.actor_id", "b.name").
	    	From("books b").
		    RunWith(bs.db)


        if filters != nil {
            if filters.Id != 0 {
                q = q.Where("id = ?", filters.Id)
            }
            if filters.Isbn != "" {
                q = q.Where("LOWER(isbn) LIKE ?", "%"+strings.ToLower(filters.Isbn)+"%") 
            }
            if filters.Title != "" {
                q = q.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(filters.Title)+"%") 
            }
            if filters.Translator != "" {
                q = q.Where("LOWER(translator) LIKE ?", "%"+strings.ToLower(filters.Translator)+"%") 
            }
            if filters.Publisher != "" {
                q = q.Where("LOWER(publisher) LIKE ?", "%"+strings.ToLower(filters.Publisher)+"%") 
            }
            if filters.Lang != "" {
                q = q.Where("LOWER(lang) LIKE ?", "%"+strings.ToLower(filters.Lang)+"%") 
            }

        }

        rows, err := q.QueryContext(ctx)
        if err != nil {
            return Result{}, err
        }

        defer rows.Close()

        var books Books
        for rows.Next() {
            var b Book
            err := rows.Scan(&b.Id, &b.Isbn, &b.Title, &b.Lang, &b.Translator, &b.Author, &b.Pages, &b.Publisher, &b.Published_date, &b.Added_date)
            if err != nil {
                return Result{}, fmt.Errorf("list characters %w,", err)
            }
            books.Response = append(books.Response, b)
        }

 

        response, err := books.Marshal()
        if err != nil {
            return Result{}, fmt.Errorf("could not marshal book %w,", err)
        }

        return Result{
            Response: response,
        }, nil
}
