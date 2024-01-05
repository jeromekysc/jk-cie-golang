package scraper

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	TARGET_SEPERATOR = "|"
)

var ALLOWED_METHODS = map[string]bool{
	"GET":  true,
	"PUT":  true,
	"POST": true,
}

var ErrMissingTarget = errors.New("Missing Target")
var ErrMethodNotAllowed = errors.New("Method not allowed")
var ErrMissingSeperator = errors.New("Target is missing seperator")

type Target struct {
	Method string
	Url    string
}

func ConvertToTarget(inc string) (*Target, error) {

	if inc == "" {
		return nil, fmt.Errorf("Missing target, %w", ErrMissingTarget)
	}

	if !strings.Contains(inc, TARGET_SEPERATOR) {
		return nil, fmt.Errorf("Invalid target, missing %s, %w", TARGET_SEPERATOR, ErrMissingSeperator)
	}

	method := strings.Split(inc, TARGET_SEPERATOR)[0]
	url := strings.Join(strings.Split(inc, TARGET_SEPERATOR)[1:], TARGET_SEPERATOR)

	if _, found := ALLOWED_METHODS[method]; !found {
		return nil, fmt.Errorf("Method %s, is not accepted %w", method, ErrMethodNotAllowed)
	}

	return &Target{
		Method: method,
		Url:    url,
	}, nil
}

func Scrape(output io.Writer, targets []Target) error {

	client := &http.Client{}

	for _, target := range targets {
		fmt.Fprintf(output, "Request : %s/%s\n", target.Method, target.Url)

		req, err := http.NewRequest(target.Method, target.Url, nil)
		resp, err := client.Do(req)

		if err != nil {
			return fmt.Errorf("error http request on %s/%s - %s", target.Method, target.Url, err)
		}

		fmt.Fprintf(output, resp.Status)
	}

	return nil
}
