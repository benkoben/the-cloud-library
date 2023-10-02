package library

// Here we will be using a SDK to communicatate with postgres
// Lets make a struct that has all the CRUD methods that full
// the dbOperator struct

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
)

// Database specifics
//
// TODO: For future improvements we could implement somekind of logic
// that determines the database driver and sets a PlaceholderFormat accordingly.
var (
    databasePlaceHolderFormat = squirrel.Dollar
)

// Errors
var (
	ErrNotFound = errors.New("not found")
)


type Book struct {
	Id             int        `json:"id"`
	Isbn           string     `json:"isbn" validate:"required"`
	Title          string     `json:"title" validate:"required"`
	Lang           string     `json:"lang" validate:"required"`
	Translator     string     `json:"translator"`
    Authors        pq.StringArray `json:"authors" validate:"required"`
	Pages          int        `json:"pages" validate:"required"`
	Publisher      string     `json:"publisher" validate:"required"`
	Published_date *time.Time `json:"published_date"`
	Added_date     *time.Time `json:"added_date"`
}

type BookStore struct {
	db *sql.DB
    placeHolderFormat squirrel.PlaceholderFormat
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
func (bs *BookStore) Store(ctx context.Context, b *Book) error {

	if b.Id == 0 {
		return bs.insert(ctx, b)
	}

	return bs.update(ctx, b)
}

// Retrieves a specific book from the database based on the id argument
func (bs *BookStore) Get(ctx context.Context, id int64) (*Book, error) {
	var b Book

    psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
    // Build query
	book := psql.Select("id, isbn, title, lang, translator, authors, pages, publisher, published_date").From("books").Where("id = ?", id).Limit(1)

    rows := book.RunWith(bs.db).QueryRowContext(ctx)
    err := rows.Scan(&b.Id, &b.Isbn, &b.Title, &b.Lang, &b.Translator, &b.Authors, &b.Pages, &b.Publisher, &b.Published_date)

    if err != nil && err != sql.ErrNoRows {
        // log the error
    }

	// Build response message
    return &b, nil
}

// Add a book to the books table
func (bs *BookStore) insert(ctx context.Context, b *Book) error {
    authors, _ := b.Authors.Value()

    // TODO: add published and added date to the insert statement
    q := squirrel.
		Insert("books").
		Columns("isbn", "title", "translator", "authors", "pages", "publisher", "lang").
		Values(b.Isbn, b.Title, b.Translator, authors, b.Pages, b.Publisher, b.Lang).
		Suffix("ON CONFLICT (isbn) DO UPDATE SET isbn = EXCLUDED.isbn, title = EXCLUDED.title, translator = EXCLUDED.translator, authors = EXCLUDED.authors, pages = EXCLUDED.pages, publisher = EXCLUDED.publisher, lang = EXCLUDED.lang").
		Suffix("RETURNING id")

    log.Println(squirrel.DebugSqlizer(q))

    return q.RunWith(bs.db).
        PlaceholderFormat(databasePlaceHolderFormat).
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
		Set("authors", pq.Array(b.Authors)).
		Set("pages", b.Pages).
		Set("publisher", b.Publisher).
		Set("published_date", b.Published_date).
		Set("added_date", b.Added_date).
		Where("id = ?", b.Id).
		RunWith(bs.db).
        PlaceholderFormat(databasePlaceHolderFormat).
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
func (bs *BookStore) Delete(ctx context.Context, b *Book) error {

	res, err := squirrel.
		Delete("books").
		Where("id = ? ", b.Id).
		RunWith(bs.db).
        PlaceholderFormat(databasePlaceHolderFormat).
		ExecContext(ctx)

	if err != nil {
		return fmt.Errorf("could not marshal book %w,", err)
	}

	rows, _ := res.RowsAffected()

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

type BooksFilters struct {
	// Id matches a books ID
	Id int
	// Isbn matches a books isbn number
	Isbn string
	// Title matches a books Title
	Title string
	// Lang matches all books published in a certain language
	Lang string
	// Translator matches all books translated by a certain Translator
	Translator string
	// Author matches all books written by a certain Author
	Author string
	// Publisher matches all books written by a certain Publisher
	Publisher string
}

// List searches for characters in the database.
//
// If filters is nil, all characters are returned. Otherwise, the results are
// filtered by the criteria in filters.
func (bs *BookStore) List(ctx context.Context, filters *BooksFilters) ([]*Book, error) {
	q := squirrel.
		Select("b.id", "b.actor_id", "b.name").
		From("books b").
		RunWith(bs.db).
        PlaceholderFormat(databasePlaceHolderFormat)

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
        // TODO: add author(s) filter here
	}

	rows, err := q.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var books []*Book
	for rows.Next() {
		var b Book
		err := rows.Scan(&b.Id, &b.Isbn, &b.Title, &b.Lang, &b.Translator, &b.Pages, &b.Publisher, &b.Published_date, &b.Added_date)
		if err != nil {
			return nil, fmt.Errorf("list characters %w,", err)
		}
		books = append(books, &b)
	}


    return books, nil
}

