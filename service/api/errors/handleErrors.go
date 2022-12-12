package errors

type ErrStatus struct {
	message string
}

func NewErrStatus(message string) error {
	errStruct := &ErrStatus{
		message: message,
	}
	return errStruct
}
func (e *ErrStatus) Error() string {
	return e.message
}
