package error

import "fmt"

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
