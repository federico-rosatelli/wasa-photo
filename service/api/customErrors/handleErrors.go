package customErrors

type StatusError struct {
	message string
}

func NewErrStatus(message string) error {
	errStruct := &StatusError{
		message: message,
	}
	return errStruct
}
func (e *StatusError) Error() string {
	return e.message
}

// This setup could be used for the costum error with the type of the error.
// e.g. an error with typeError equal to 0 would be an internalstatuserror.
// The function Type is used like a switch case for the type error so you can delete
// the const variables in the file routing-functions and use only the int value
// type StatusError struct {
// 	message   string
// 	typeError int
// }

// type error interface {
// 	Error() string
// 	Type() int
// }

// func NewErrStatus(message string) error {
// 	errStruct := &StatusError{
// 		message: message,
// 	}
// 	return errStruct
// }
// func (e *StatusError) Error() string {
// 	return e.message
// }

// func (e *StatusError) Type() int {
// 	return e.typeError
// }
