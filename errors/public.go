package errors

// Public creates an error that includes a user-friendly message
// while wrapping the underlying error. This allows you to display
// a custom message to users and still retain the original error 
// details for internal debugging.
func Public(err error, msg string) error {
	return publicError{err, msg}
}

type publicError struct {
	err error
	msg string
}

func (pe publicError) Error() string {
	return pe.err.Error()
}

func (pe publicError) Public() string {
	return pe.msg
}

func (pe publicError) Unwrap() error {
	return pe.err
}
