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

package readerresult

import (
	"context"
	"errors"
	"fmt"

	ER "github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/result"
)

// ExampleMapLeft adds context to the error channel. The original error stays
// in the chain, so errors.Is keeps working.
func ExampleMapLeft() {
	errNotFound := errors.New("not found")

	loadConfig := func(name string) ReaderResult[string] {
		return F.Pipe1(
			Left[string](errNotFound),
			MapLeft[string](ER.OnError("loading config %q", name)),
		)
	}

	_, err := result.Unwrap(loadConfig("app.yaml")(context.Background()))
	fmt.Println(err)
	fmt.Println(errors.Is(err, errNotFound))

	// Output:
	// loading config "app.yaml", Caused By: not found
	// true
}
