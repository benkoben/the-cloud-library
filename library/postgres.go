package library

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

var (
	postgresBaseConnString = "postgres://{username}:{password}@{host}/{database}?sslmode={sslmode}"
)


type Credentials struct {
	Username string
	Password string
}

type postgresClient struct {
	connString string
    context context.Context
    driver string
}

func (client postgresClient) InitDb() (*sql.DB, error) {
	db, err := sql.Open(client.driver, client.connString)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}

type PostgresClientOptions struct {
	Host       string
	SslEnabled bool
	Database   string
}


func NewPgClient(cred Credentials, options PostgresClientOptions) (*sql.DB, error) {
	if options.Database == "" {
		return nil, errors.New("database must not be empty")
	}

	replacer := strings.NewReplacer(
		"{username}", cred.Username,
		"{password}", cred.Password,
		"{host}", options.Host,
		"{database}", options.Database,
		"{sslmode}", strconv.FormatBool(options.SslEnabled),
	)
	connStr := replacer.Replace(postgresBaseConnString)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

    return db, nil
}
