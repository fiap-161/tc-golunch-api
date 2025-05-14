package errors

type ValidationError struct {
	Msg string
}

func (e *ValidationError) Error() string {
	return e.Msg
}

type UnauthorizedError struct {
	Msg string
}

func (e *UnauthorizedError) Error() string {
	return e.Msg
}
