package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

const (
    defaultTimeout = 10
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

	fileContent, err := os.ReadFile(path)
    
    if len(fileContent) == 0 {
        return nil, fmt.Errorf("file not found. make sure that %s exists and containts valid db credentials", path)
    }

	if err != nil {
		return nil, fmt.Errorf("could not read file in path %s: %s\n", path, err)
	}

	var cred Credentials
	if err := json.Unmarshal(fileContent, &cred); err != nil {
		return nil, fmt.Errorf("could not unmarshal %s into credentials: %s\n",path, err)
	}
	return &cred, nil
}

type PostgresClient struct {
    Client     *sql.DB
    Timeout    time.Duration
}

type PostgresClientOptions struct {
	Host       string
	SslEnabled bool
	Database   string
    Timeout    time.Duration
}

func NewPgClient(cred Credentials, options PostgresClientOptions) (*PostgresClient, error) {
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

	return &PostgresClient{
        Client: db,
        Timeout: options.Timeout,
    }, db.Ping()
}

// Checks if the database connection is alive. If possible it will restart the connection.
func (c *PostgresClient) IsHealthy(ctx context.Context, healthCh chan <- struct{}){
    if err := c.Client.PingContext(ctx); err != nil {
                 
    }
    return
}