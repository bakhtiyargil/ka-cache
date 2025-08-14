package http

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	InternalServerError   = errors.New("internal server error")
	ResourceNotFoundError = errors.New("resource not found error")
)

type RestError interface {
	Status() int
	Error() string
	Causes() interface{}
}

type RestErrorResponse struct {
	ErrStatus int         `json:"status,omitempty"`
	ErrError  string      `json:"error,omitempty"`
	ErrCauses interface{} `json:"-"`
}

func NewResourceNotFound(causes interface{}) RestError {
	result := &RestErrorResponse{
		ErrStatus: http.StatusNotFound,
		ErrError:  ResourceNotFoundError.Error(),
		ErrCauses: causes,
	}
	return result
}

func NewInternalServerError(causes interface{}) RestError {
	result := &RestErrorResponse{
		ErrStatus: http.StatusInternalServerError,
		ErrError:  InternalServerError.Error(),
		ErrCauses: causes,
	}
	return result
}

func (e *RestErrorResponse) Error() string {
	return fmt.Sprintf("status: %d - errors: %s - causes: %v", e.ErrStatus, e.ErrError, e.ErrCauses)
}

func (e *RestErrorResponse) Status() int {
	return e.ErrStatus
}

func (e *RestErrorResponse) Causes() interface{} {
	return e.ErrCauses
}
