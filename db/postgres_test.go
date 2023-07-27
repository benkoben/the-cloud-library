package db

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewPgOperator(t *testing.T) {
	var tests = []struct {
		name      string
		cred      Credentials
		options   PostgresClientOptions
		want      *postgresClient
		wantError error
	}{
		{
			name: "new client",
			options: PostgresClientOptions{
				Host:       "localhost",
				SslEnabled: false,
				Database:   "library",
			},
			want: &postgresClient{
				connString: "postgres://postgres:Syp9393@localhost/library?sslmode=false",
				driver:     "postgres",
			},
			cred: Credentials{
				Username: "postgres",
				Password: "Syp9393",
			},
		},
	}

	for _, test := range tests {

		got, gotErr := NewPgClient(test.cred, test.options)

		if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(postgresClient{})); diff != "" {
			t.Errorf("NewPgClient(%+v, %+v) = unexpected results, (-want, +got)\n%s\n", test.cred, test.options, diff)
		}

		if test.wantError != nil && gotErr == nil {
			t.Errorf("Unexpected result, should return error")
		}
	}
}
