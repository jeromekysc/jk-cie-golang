package scraper

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToTarget(t *testing.T) {

	tests := []struct {
		input    string
		expected *Target
	}{
		{
			input: "GET|https:/www.google.com",
			expected: &Target{
				Method: "GET",
				Url:    "https:/www.google.com",
			},
		},
		{
			input:    "DELETE|https:/www.google.com",
			expected: nil,
		},
		{
			input:    "https:/www.google.com",
			expected: nil,
		},
		{
			input:    "",
			expected: nil,
		},
	}

	for _, tc := range tests {

		t.Run("Testing "+tc.input, func(t *testing.T) {
			target, err := ConvertToTarget(tc.input)

			if tc.expected == nil {
				assert.NotNil(t, err)
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
