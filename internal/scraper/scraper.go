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

type CanDoRollback interface {
	CanDoRollback() bool
	DoRollback() error
}

type Rollbackable interface {
	Rollback() error
}

type ErrMissingTargetType struct {
	OriginalMethod string
	InnerStruct    Rollbackable
}

func (e *ErrMissingTargetType) GetAllowedMethods() map[string]bool {
	return ALLOWED_METHODS
}

func (e *ErrMissingTargetType) GetOriginalMethod() string {
	return e.OriginalMethod
}

func (e *ErrMissingTargetType) Error() string {
	return ErrMissingTarget.Error()
}

func (e *ErrMissingTargetType) Is(target error) bool {
	t, ok := target.(*ErrMissingTargetType)

	if !ok {
		return false
	}

	return t.GetOriginalMethod() == e.GetOriginalMethod()
}

func (e *ErrMissingTargetType) CanDoRollback() bool {
	return true
}

func (e *ErrMissingTargetType) DoRollback() error {
	return e.InnerStruct.Rollback()
}

type Target struct {
	Method string
	Url    string
}

type User struct {
	name string
}

func (u *User) Rollback() error {
	return nil
}

func ConvertToTarget(inc string) (*Target, error) {

	obj := User{"miguel"}

	if inc == "" {
		return nil, &ErrMissingTargetType{
			InnerStruct: &obj,
		}
	}

	if !strings.Contains(inc, TARGET_SEPERATOR) {
		return nil, fmt.Errorf("invalid target, missing %s, %w", TARGET_SEPERATOR, ErrMissingSeperator)
	}

	method := strings.Split(inc, TARGET_SEPERATOR)[0]
	url := strings.Join(strings.Split(inc, TARGET_SEPERATOR)[1:], TARGET_SEPERATOR)

	if _, found := ALLOWED_METHODS[method]; !found {
		return nil, fmt.Errorf("method %s, is not accepted %w", method, ErrMethodNotAllowed)
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
