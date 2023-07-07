package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnError(t *testing.T) {
	onError := OnError("Failed to process [%s]", "filename")

	err := fmt.Errorf("Cause")

	err1 := onError(err)

	assert.Equal(t, "Failed to process [filename], Caused By: Cause", err1.Error())
}
