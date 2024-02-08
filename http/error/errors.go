package error

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	InternalServerError = errors.New("internal Server Error")
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

func NewInternalServerError(causes interface{}) RestError {
	result := RestErrorResponse{
		ErrStatus: http.StatusInternalServerError,
		ErrError:  InternalServerError.Error(),
		ErrCauses: causes,
	}
	return result
}

func ParseErrors(err error) RestError {
	switch {
	default:
		if restErr, ok := err.(RestError); ok {
			return restErr
		}
		return NewInternalServerError(err)
	}
}

func (e RestErrorResponse) Error() string {
	return fmt.Sprintf("status: %d - errors: %s - causes: %v", e.ErrStatus, e.ErrError, e.ErrCauses)
}

func (e RestErrorResponse) Status() int {
	return e.ErrStatus
}

func (e RestErrorResponse) Causes() interface{} {
	return e.ErrCauses
}

func NewRestErrorResponse(status int, err string, causes interface{}) RestError {
	return RestErrorResponse{
		ErrStatus: status,
		ErrError:  err,
		ErrCauses: causes}
}

func ErrorResponse(err error) (int, interface{}) {
	return ParseErrors(err).Status(), ParseErrors(err)
}
