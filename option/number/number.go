//   Copyright (c) 2023 IBM Corp.
//   All rights reserved.
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package number

import (
	"strconv"

	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
)

func atoi(value string) (int, bool) {
	data, err := strconv.Atoi(value)
	return data, err == nil
}

var (
	// Atoi converts a string to an integer
	Atoi = O.Optionize1(atoi)
	// Itoa converts an integer to a string
	Itoa = F.Flow2(strconv.Itoa, O.Of[string])
)
