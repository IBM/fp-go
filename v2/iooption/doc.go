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

// Package iooption provides the IOOption monad, combining IO effects with Option for optional values.
//
// # Fantasy Land Specification
//
// This is a monad transformer combining:
//   - IO monad: https://github.com/fantasyland/fantasy-land
//   - Maybe (Option) monad: https://github.com/fantasyland/fantasy-land#maybe
//
// Implemented Fantasy Land algebras:
//   - Functor: https://github.com/fantasyland/fantasy-land#functor
//   - Apply: https://github.com/fantasyland/fantasy-land#apply
//   - Applicative: https://github.com/fantasyland/fantasy-land#applicative
//   - Chain: https://github.com/fantasyland/fantasy-land#chain
//   - Monad: https://github.com/fantasyland/fantasy-land#monad
//   - Alt: https://github.com/fantasyland/fantasy-land#alt
//   - Plus: https://github.com/fantasyland/fantasy-land#plus
//   - Alternative: https://github.com/fantasyland/fantasy-land#alternative
//
// IOOption[A] represents a computation that:
//   - Performs side effects (IO)
//   - May or may not produce a value of type A (Option)
//
// This is defined as: IO[Option[A]] or func() Option[A]
package iooption

//go:generate go run .. iooption --count 10 --filename gen.go
