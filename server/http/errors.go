package http

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	InternalServerError   = errors.New("internal_server_error")
	ResourceNotFoundError = errors.New("resource_not_found_error")
)

type RestError interface {
	RequestId() string
	Status() int
	Error() string
	Causes() interface{}
}

type RestErrorResponse struct {
	ErrRequestId string      `json:"requestId,omitempty"`
	ErrStatus    int         `json:"status,omitempty"`
	ErrError     string      `json:"error,omitempty"`
	ErrCauses    interface{} `json:"causes,omitempty"`
}

func NewResourceNotFound(requestId string, causes interface{}) RestError {
	result := &RestErrorResponse{
		ErrRequestId: requestId,
		ErrStatus:    http.StatusNotFound,
		ErrError:     ResourceNotFoundError.Error(),
		ErrCauses:    causes,
	}
	return result
}

func NewInternalServerError(requestId string, causes interface{}) RestError {
	result := &RestErrorResponse{
		ErrRequestId: requestId,
		ErrStatus:    http.StatusInternalServerError,
		ErrError:     InternalServerError.Error(),
		ErrCauses:    causes,
	}
	return result
}

func (e *RestErrorResponse) RequestId() string {
	return e.ErrRequestId
}

func (e *RestErrorResponse) Status() int {
	return e.ErrStatus
}

func (e *RestErrorResponse) Error() string {
	return fmt.Sprintf("status: %d - errors: %s - causes: %v", e.ErrStatus, e.ErrError, e.ErrCauses)
}

func (e *RestErrorResponse) Causes() interface{} {
	return e.ErrCauses
}
