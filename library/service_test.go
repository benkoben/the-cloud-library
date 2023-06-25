package library

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestNewService(t *testing.T) {
	var tests = []struct {
		name          string
		inputOperator database
		inputOptions  ServiceOptions
		want          *service
		wantError     error
	}{
		{
			name:          "new service",
			inputOperator: fakeDb{},
			inputOptions:  ServiceOptions{},
			want: &service{
				database:   fakeDb{},
				timeout:    time.Second * 10,
				concurreny: 1,
			},
			wantError: nil,
		},
	}

	for _, test := range tests {
		got, gotErr := NewService(test.inputOperator, test.inputOptions)

		if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(service{}, fakeDb{})); diff != "" {
			t.Errorf("NewService(%q, %q) = unexpected results, (-want, +got)\n%s\n", test.inputOperator, test.inputOptions, diff)
		}

		if test.wantError != nil && gotErr == nil {
			t.Errorf("Unexpected result, should return error")
		}
	}
}

type fakeDb struct {
    connErr bool   
}

func (db fakeDb) InitDb() error {
    if db.connErr {
        return errors.New("cannot connect to db")
    }
    return nil
} 
