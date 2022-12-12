package errors

import "fmt"

type ErrStatus struct {
	message string
}

func NewErrStatus(message string) error {
	errStruct := &ErrStatus{
		message: message,
	}
	errRet := fmt.Errorf("%w", errStruct)
	return errRet
}
func (e *ErrStatus) Error() string {
	return e.message
}
