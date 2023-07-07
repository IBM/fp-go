package utils

import (
	"errors"
	"strings"
)

var Upper = strings.ToUpper

func Sum(left, right int) int {
	return left + right
}

func Double(value int) int {
	return value * 2
}

func Triple(value int) int {
	return value * 3
}

func StringLen(value string) int {
	return len(value)
}

func Error() (int, error) {
	return 0, errors.New("some error")
}
