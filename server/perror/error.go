package perror

const internalServerError = "Internal error, please check logs"

type PError struct {
	ErrorMessage string
	err          error
}

func (e *PError) Error() error {
	if e.err != nil {
		return e.err
	}
	return nil
}

func (e *PError) String() string {
	return e.Error().Error()
}

func (e *PError) Message() string {
	return e.ErrorMessage
}

func (e *PError) AsJSON() string {
	return `{"error": "` + e.Message() + `"}`
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
