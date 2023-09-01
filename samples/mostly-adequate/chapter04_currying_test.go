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

package mostlyadequate

import (
	"regexp"
	"strings"

	F "github.com/IBM/fp-go/function"
	N "github.com/IBM/fp-go/number"
	I "github.com/IBM/fp-go/number/integer"
	S "github.com/IBM/fp-go/string"
)

var (
	Match = F.Curry2((*regexp.Regexp).FindStringSubmatch)
	Split = F.Curry2(F.Bind3of3((*regexp.Regexp).Split)(-1))

	Add      = N.Add[int]
	ToString = I.ToString
	ToLower  = strings.ToLower
	ToUpper  = strings.ToUpper
	Concat   = F.Curry2(S.Monoid.Concat)
)

// Replace cannot be generated via [F.Curry3] because the order of parameters does not match our desired curried order
func Replace(search *regexp.Regexp) func(replace string) func(s string) string {
	return func(replace string) func(s string) string {
		return func(s string) string {
			return search.ReplaceAllString(s, replace)
		}
	}
}
