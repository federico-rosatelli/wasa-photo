package errors

type ErrStatus struct {
	message string
}

func NewErrStatus(message string) *ErrStatus {
	return &ErrStatus{
		message: message,
	}
}
func (e *ErrStatus) Error() string {
	return e.message
}
