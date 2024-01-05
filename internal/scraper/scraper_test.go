package scraper

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToTarget(t *testing.T) {

	tests := []struct {
		input     string
		expected  *Target
		expectErr error
	}{
		{
			input: "GET|https:/www.google.com",
			expected: &Target{
				Method: "GET",
				Url:    "https:/www.google.com",
			},
		},
		// {
		// 	input:     "FOOBAR|https:/www.google.com",
		// 	expected:  nil,
		// 	expectErr: ErrMethodNotAllowed,
		// },
		// {
		// 	input:     "https:/www.google.com",
		// 	expected:  nil,
		// 	expectErr: ErrMissingSeperator,
		// },
		{
			input:     "",
			expected:  nil,
			expectErr: &ErrMissingTargetType{},
		},
	}

	for _, tc := range tests {

		t.Run("Testing "+tc.input, func(t *testing.T) {
			target, err := ConvertToTarget(tc.input)

			assert.ErrorIs(t, err, tc.expectErr)

			if err != nil {
				if castedErr, worked := err.(*ErrMissingTargetType); worked {
					fmt.Println(castedErr.GetOriginalMethod())
				}

				switch t := err.(type) {
				case CanDoRollback:
					if t.CanDoRollback() {
						t.Rollback() // you don't care which type it is
					}

					fallthrough
				case *ErrMissingTargetType:
					if t.
					fmt.Println(t.GetOriginalMethod())
				default:
					fmt.Println("different err")
				}

			}

			assert.Equal(
				t,
				target,
				tc.expected,
			)
		})
	}

}

func TestScrape(t *testing.T) {
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	outBuffer := bytes.Buffer{}

	tests := []Target{
		{
			Method: "GET",
			Url:    fakeServer.URL,
		},
		{
			Method: "PUT",
			Url:    fakeServer.URL,
		},
		{
			Method: "POST",
			Url:    fakeServer.URL,
		},
	}

	for _, tc := range tests {

		t.Run("Testing "+tc.Method+"/"+tc.Url, func(t *testing.T) {

			err := Scrape(
				&outBuffer,
				[]Target{
					{
						Method: tc.Method,
						Url:    tc.Url,
					},
				},
			)

			if err != nil {
				t.Error(err)
			}

			assert.Contains(
				t,
				outBuffer.String(),
				tc.Method,
			)

			assert.Contains(
				t,
				outBuffer.String(),
				tc.Url,
			)
		})

	}

}
