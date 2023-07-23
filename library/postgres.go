package library

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var (
	postgresBaseConnString = "postgres://{username}:{password}@{host}/{database}?sslmode={sslmode}"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Credentials are always read from filesystem. This way we can mount our credentials to our docker container/pod
// Which is generally more secure than using environment variables for credentials.
func NewDbCredentials(path string) (*Credentials, error) {

	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Errorf("could not read file: %s\n", err)
	}
	var cred Credentials
	if err := json.Unmarshal(fileContent, &cred); err != nil {
		return nil, fmt.Errorf("could not unmarshal file into credentials: %s\n", err)
	}
	return &cred, nil
}

type postgresClient struct {
	connString string
	context    context.Context
	driver     string
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
