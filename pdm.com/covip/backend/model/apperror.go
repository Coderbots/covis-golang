package model

import "net/http"

type APIError interface {
	// APIError returns HTTP statuscode and error message
	APIError() (int, string)
}

type sentinelAPIError struct {
	status int
	msg    string
}

func (e sentinelAPIError) Error() string {
	return e.msg
}

func (e sentinelAPIError) APIError() (int, string) {
	return e.status, e.msg
}

var (
	ErrNotFound = &sentinelAPIError{status: http.StatusNotFound, msg: "not found"}
)

type sentinelWrappedError struct {
	error
	sentinel *sentinelAPIError
}

func (e sentinelWrappedError) Is(err error) bool {
	return e.sentinel == err
}

func (e sentinelWrappedError) APIError() (int, string) {
	return e.sentinel.APIError()
}

func WrapError(err error, sentinel *sentinelAPIError) error {
	return sentinelWrappedError{error: err, sentinel: sentinel}
}
