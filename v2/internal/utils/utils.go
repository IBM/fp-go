// Copyright (c) 2023 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"errors"
	"strings"
)

var Upper = strings.ToUpper

func Inc(i int) int {
	return i + 1
}

func Dec(i int) int {
	return i - 1
}

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
