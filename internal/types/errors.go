package types

import (
	"errors"
	"net/http"
)

// @name HTTPError
// @description Standard error response returned by the API
// @property Code
// @type integer
// @description HTTP status code
// @property Err
// @type error
// @description Human-readable error message
type HTTPErrorSwaggerWrapper HTTPError

// HTTPError represents an error response returned by the API.
// This is used to standardize error responses across the application.
type HTTPError struct {
	// HTTP status code
	Code int `json:"code"`

	// Error message
	Err error `json:"error"`
}

func (e HTTPError) Error() string {
	return e.Err.Error()
}

var (
	ErrBadRequest          = func(err error) HTTPError { return HTTPError{Code: http.StatusBadRequest, Err: err} }
	ErrNotFound            = func(err error) HTTPError { return HTTPError{Code: http.StatusNotFound, Err: err} }
	ErrInternalServerError = func(err error) HTTPError { return HTTPError{Code: http.StatusInternalServerError, Err: err} }
	ErrConflict            = func(err error) HTTPError { return HTTPError{Code: http.StatusConflict, Err: err} }
)

var (
	ErrWalletNotFound    = errors.New("wallet not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrConcurrentUpdate  = errors.New("concurrent update detected, retry required")
	ErrInvalidOperation  = errors.New("invalid operation type, must be DEPOSIT or WITHDRAW")
	ErrOperationExists   = errors.New("operation with this reference_id already exists")
	ErrDB                = errors.New("database error")
)
