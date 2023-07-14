package io

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {

	l := Logger[int]()

	lio := l("out")

	assert.Equal(t, nil, lio(10)())
}

func TestLogf(t *testing.T) {

	l := Logf[int]()

	lio := l("Value is %d")

	assert.Equal(t, nil, lio(10)())
}
