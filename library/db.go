package library

// Here we will be using a SDK to communicatate with postgres
// Lets make a struct that has all the CRUD methods that fullfill
// the dbOperator struct

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
    "github.com/go-playground/validator/v10"
)

var (
	book_insertion_query = `
        INSERT INTO books (
            id, isbn,title,
            lang, translator,
            author, pages, publisher, 
            published_date,added_date,
        ) VALUES (, {title},
                  {lang}, {translator},
                  {author}, {pages},
                  {publisher},{published_date},
                  {added_date}
        )
    `
)

type dbClient interface {
    ExecContext(context.Context, string, ...any) (sql.Result, error)
    Close() error
}

type Credentials struct {
	Username string
	Password string
}

type bookData struct {
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

func (b bookData) validate() error {
    validate := validator.New()
    err := validate.Struct(b)
    if err != nil {
        return errors.New("payload is missing required keys")

    }
    return nil
}


type Book struct {
    Client dbClient
	Data bookData
}

func (b Book)templateBookQuery(template string) string {
	replacer := strings.NewReplacer(
		"{id}", strconv.Itoa(b.Data.Id),
		"{isbn}", b.Data.Isbn,
		"{title}", b.Data.Title,
		"{lang}", b.Data.Lang,
		"{translator}", b.Data.Translator,
		"{author}", b.Data.Author,
		"{pages}", strconv.Itoa(b.Data.Pages),
		"{publisher}", b.Data.Publisher,
		"{published_date}", b.Data.Published_date.String(),
		"{added_date}", b.Data.Added_date.String(),
	)
    return replacer.Replace(template)
}

// Adds a book to database
func (b Book) Add(ctx context.Context) (Result, error) {
	query := b.templateBookQuery(book_insertion_query)
	result, err := b.Client.ExecContext(ctx, query)
	if err != nil {
		return Result{}, err
	}

	return Result{
        table: "book",
        result: result,
    }, nil
}


//func (b book) Remove() (sql.Result, error) {}
//func (b book) Update() (sql.Result, error) {}

// constructor method that is used to unmarshall a REST request into a book object
// which has CRUD methods.
func NewBook(data []byte) (*Book, error) {
	var bookData *bookData

    if ok := json.Valid(data); !ok {
        return nil, errors.New("payload is not in a valid json format")
    }

	err := json.Unmarshal(data, &bookData)
	if err != nil {
		return nil, err
	}
   
    if err := bookData.validate(); err != nil {
        return nil, err
    }

	return &book{
        Data: *bookData,
    }, nil
}

