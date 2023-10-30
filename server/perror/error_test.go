package perror

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPError(t *testing.T) {
	err := errors.New("test error")
	message := "test message"
	pErr := NewPError(err, message)

	assert.Equal(t, `{"error": "`+message+`"}`, pErr.AsJSON())
	assert.Equal(t, err.Error(), pErr.Error())
	assert.Equal(t, err, pErr.Err())
	assert.Equal(t, message, pErr.Message())
}
