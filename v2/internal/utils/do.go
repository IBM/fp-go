// Copyright (c) 2023 - 2025 IBM Corp.
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
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
)

type (
	Initial struct {
	}

	WithLastName struct {
		Initial
		LastName string
	}

	WithGivenName struct {
		WithLastName
		GivenName string
	}
)

var (
	Empty = Initial{}

	SetLastName = F.Curry2(func(name string, s1 Initial) WithLastName {
		return WithLastName{
			Initial:  s1,
			LastName: name,
		}
	})

	SetGivenName = F.Curry2(func(name string, s1 WithLastName) WithGivenName {
		return WithGivenName{
			WithLastName: s1,
			GivenName:    name,
		}
	})
)

func GetFullName(s WithGivenName) string {
	return fmt.Sprintf("%s %s", s.GivenName, s.LastName)
}
