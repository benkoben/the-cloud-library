package library

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestNewBook(t *testing.T) {
	var tests = []struct {
		name      string
		input     string
		want      *book
		wantError error
	}{
		{
			name:  "new book with valid payload",
			input: jsonData,
			want: &book{
				Data: bookData{
					Id:             1234,
					Isbn:           "9789100187934",
					Title:          "Pesten",
					Lang:           "swedish",
					Translator:     "Jan Stolpe",
					Author:         "Albert Camus",
					Pages:          254,
					Publisher:      "Albert Bonniers",
					Published_date: stringToTime("2022-03-02"),
					Added_date:     stringToTime("2022-03-02"),
				},
			},
			wantError: nil,
		},
		{
			name:      "new book with invalid payload",
			input:     errJsonData,
			want:      nil,
			wantError: errors.New("error"),
		},
	}

	for _, test := range tests {
		got, gotErr := NewBook([]byte(test.input))

		if diff := cmp.Diff(test.want, got); diff != "" {
			t.Errorf("NewBook(%q) = unexpected results, (-want, +got)\n%s\n", test.input, diff)
		}

		if test.wantError != nil && gotErr == nil {
			t.Errorf("Unexpected result, should return error")
		}
	}
}

var jsonData string = `{
			"id": 1234,
			"isbn": "9789100187934",
			"title":"Pesten",
			"lang":"swedish",
			"translator":"Jan Stolpe",
			"author":"Albert Camus",
			"pages": 254,
			"publisher":"Albert Bonniers",
			"published_date":"2022-03-02T00:00:00Z",
			"added_date":"2022-03-02T00:00:00Z"
		}
	`
var errJsonData string = `
    {
        "can_this_be_marshalled": false        
    }
`

func stringToTime(sTime string) *time.Time {
	t, err := time.Parse(time.DateOnly, sTime)
	if err != nil {
		fmt.Println(err)
	}
	return &t
}
