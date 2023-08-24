package readiness

type UnrecoverableError struct {
	desc string
	err  error
}

func (e *UnrecoverableError) Error() string {
	return e.desc + ": " + e.err.Error()
}

func (e *UnrecoverableError) Unwrap() error {
	return e.err
}

// NewUnrecoverableError creates an Error object that is not recoverable
func NewUnrecoverableError(err error) error {
	return &UnrecoverableError{
		desc: "Unrecoverable",
		err:  err,
	}
}
