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

package match

import (
	"fmt"

	"github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/result"
	S "github.com/IBM/fp-go/v2/string"
)

type Thing struct {
	Name string
}

func (t Thing) GetName() string {
	return t.Name
}

var (
	// func(Thing) Either[error, string]
	getName = F.Flow2(
		Thing.GetName,
		result.FromPredicate(S.IsNonEmpty, errors.OnSome[string]("value [%s] is empty")),
	)

	// func(option.Option[Thing]) Either[error, string]
	GetName = F.Flow2(
		result.FromOption[Thing](errors.OnNone("value is none")),
		result.Chain(getName),
	)
)

func ExampleEither_match() {

	oThing := O.Of(Thing{"Carsten"})

	res := F.Pipe2(
		oThing,
		GetName,
		result.Fold(S.Format[error]("failed with error %v"), S.Format[string]("get value %s")),
	)

	fmt.Println(res)

	// Output:
	// get value Carsten

}
