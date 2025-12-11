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

// Package ioresult provides the IOResult monad, combining IO effects with Result for error handling.
//
// # Fantasy Land Specification
//
// This is a monad transformer combining:
//   - IO monad: https://github.com/fantasyland/fantasy-land
//   - Either monad: https://github.com/fantasyland/fantasy-land#either
//
// Implemented Fantasy Land algebras:
//   - Functor: https://github.com/fantasyland/fantasy-land#functor
//   - Bifunctor: https://github.com/fantasyland/fantasy-land#bifunctor
//   - Apply: https://github.com/fantasyland/fantasy-land#apply
//   - Applicative: https://github.com/fantasyland/fantasy-land#applicative
//   - Chain: https://github.com/fantasyland/fantasy-land#chain
//   - Monad: https://github.com/fantasyland/fantasy-land#monad
//   - Alt: https://github.com/fantasyland/fantasy-land#alt
//
// IOResult[A] represents a computation that:
//   - Performs side effects (IO)
//   - Can fail with an error or succeed with a value of type A (Result/Either)
//
// This is defined as: IO[Result[A]] or func() Either[error, A]
// IOResult is an alias for IOEither with error as the left type.
package ioresult

//go:generate go run .. ioeither --count 10 --filename gen.go
