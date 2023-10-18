package perror

const internalServerError = "Internal error, please check logs"

type PError struct {
	ErrorMessage string `json:"error_message"`
	err          error
}

func (e *PError) Error() string {
	return e.err.Error()
}

func (e *PError) String() string {
	return e.Error()
}

func (e *PError) Message() string {
	return e.ErrorMessage
}

func NewPError(err error, message string) *PError {
	return &PError{
		ErrorMessage: message,
		err:          err,
	}
}

func NewSinglePError(err error) *PError {
	return NewPError(err, err.Error())
}

func NewInternalServerPError(err error) *PError {
	return NewPError(err, internalServerError)
}
