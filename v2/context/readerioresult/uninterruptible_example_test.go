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

package readerioresult

import (
	"context"
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
)

// ExampleUninterruptible shows that an uninterruptible computation runs even if the
// context has already been cancelled.
func ExampleUninterruptible() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	flush := F.Pipe1(
		FromIO(func() string { return "flushed" }),
		WithContext[string],
	)

	fmt.Println(flush(ctx)())
	fmt.Println(Uninterruptible(flush)(ctx)())
	// Output:
	// Left[*errors.errorString](context canceled)
	// Right[string](flushed)
}
