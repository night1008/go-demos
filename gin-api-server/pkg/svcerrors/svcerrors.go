package svcerrors

import (
	"errors"
	"fmt"
)

type SvcError struct {
	Code    string         `json:"code"`    // service error code
	Status  int            `json:"-"`       // http status code
	Message string         `json:"message"` // descriptive message for caller
	Errs    []SvcErrorItem `json:"errors"`  // concrete error content for server
}

type SvcErrorItem struct {
	Key string `json:"key"`   // form path if it is a validation error, otherwise just error index
	Err string `json:"error"` // concrete error content
}

var codes = map[string]string{}

func NewError(code string, status int, message string) *SvcError {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("error code %s already exist", code))
	}
	codes[code] = message
	return &SvcError{Code: code, Status: status, Message: message, Errs: []SvcErrorItem{}}
}

func (e *SvcError) Error() string {
	return fmt.Sprintf("code: %s, message: %s", e.Code, e.Message)
}

func (e *SvcErrorItem) Error() string {
	return e.Err
}

func (e *SvcError) WithErrors(errs ...error) *SvcError {
	newError := *e
	newError.Errs = []SvcErrorItem{}
	for i, err := range errs {
		var item *SvcErrorItem
		if errors.As(err, &item) {
			newError.Errs = append(newError.Errs, *item)
		} else {
			newError.Errs = append(newError.Errs, SvcErrorItem{Key: fmt.Sprintf("%d", i), Err: err.Error()})
		}
	}
	return &newError
}
