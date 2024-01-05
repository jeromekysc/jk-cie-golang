package scraper

import (
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
)

const (
	TARGET_SEPERATOR = "|"
)

var ALLOWED_METHODS = []string{"GET", "POST", "PUT"}

type Target struct {
	Method string
	Url    string
}

func ConvertToTarget(inc string) (*Target, error) {

	if inc == "" {
		return nil, fmt.Errorf("Missing target")
	}

	if !strings.Contains(inc, TARGET_SEPERATOR) {
		return nil, fmt.Errorf("Invalid target, missing " + TARGET_SEPERATOR)
	}

	method := strings.Split(inc, TARGET_SEPERATOR)[0]
	url := strings.Join(strings.Split(inc, TARGET_SEPERATOR)[1:], TARGET_SEPERATOR)

	if !slices.Contains(ALLOWED_METHODS, method) {
		return nil, fmt.Errorf("Method %s, is not accepted", method)
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
