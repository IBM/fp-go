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

package readeroption

import (
	"github.com/IBM/fp-go/v2/internal/readert"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
)

// Sequence swaps the order of nested environment parameters in a ReaderOption computation.
//
// This function takes a ReaderOption that produces another ReaderOption and returns a
// reader.Kleisli that reverses the order of the environment parameters. The result is
// a curried function that takes R2 first, then R1, and produces an Option[A].
//
// Type Parameters:
//   - R1: The first environment type (becomes inner after flip)
//   - R2: The second environment type (becomes outer after flip)
//   - A: The value type
//
// Parameters:
//   - ma: A ReaderOption that takes R2 and may produce a ReaderOption[R1, A]
//
// Returns:
//   - A reader.Kleisli[R2, R1, Option[A]], which is func(R2) func(R1) Option[A]
//
// The function preserves optional values at both levels. If either the outer or inner
// computation produces None, the final result will be None.
//
// Example:
//
//	import S "github.com/IBM/fp-go/v2/string"
//
//	type Database struct {
//	    ConnectionString string
//	}
//	type Config struct {
//	    Timeout int
//	}
//
//	// Original: takes Config, may produce ReaderOption[Database, string]
//	original := func(cfg Config) option.Option[ReaderOption[Database, string]] {
//	    if cfg.Timeout <= 0 {
//	        return option.None[ReaderOption[Database, string]]()
//	    }
//	    return option.Some(func(db Database) option.Option[string] {
//	        if S.IsEmpty(db.ConnectionString) {
//	            return option.None[string]()
//	        }
//	        return option.Some(fmt.Sprintf("Query on %s with timeout %d",
//	            db.ConnectionString, cfg.Timeout))
//	    })
//	}
//
//	// Sequenced: takes Database first, then Config
//	sequenced := Sequence(original)
//
//	db := Database{ConnectionString: "localhost:5432"}
//	cfg := Config{Timeout: 30}
//
//	// Apply database first to get a function that takes config
//	configReader := sequenced(db)
//	// Then apply config to get the final result
//	result := configReader(cfg)
//	// result is Option[string]
func Sequence[R1, R2, A any](ma ReaderOption[R2, ReaderOption[R1, A]]) reader.Kleisli[R2, R1, Option[A]] {
	return readert.Sequence(
		option.Chain,
		ma,
	)
}

// SequenceReader swaps the order of environment parameters when the inner computation is a Reader.
//
// This function is similar to Sequence but specialized for the case where the innermost computation
// is a pure Reader (without optional values) rather than another ReaderOption. It takes a
// ReaderOption that produces a Reader and returns a Reader that produces a ReaderOption.
//
// Type Parameters:
//   - R1: The first environment type (becomes outer after flip)
//   - R2: The second environment type (becomes inner after flip)
//   - A: The value type
//
// Parameters:
//   - ma: A ReaderOption that takes R2 and may produce a Reader[R1, A]
//
// Returns:
//   - A reader.Kleisli[R2, R1, Option[A]], which is func(R2) func(R1) Option[A]
//
// The function preserves the optional nature from the outer ReaderOption layer. If the outer
// computation produces None, the inner ReaderOption result will be None.
//
// Example:
//
//	type Database struct {
//	    ConnectionString string
//	}
//	type Config struct {
//	    Timeout int
//	}
//
//	// Original: takes Config, may produce Reader[Database, string]
//	original := func(cfg Config) option.Option[Reader[Database, string]] {
//	    if cfg.Timeout <= 0 {
//	        return option.None[Reader[Database, string]]()
//	    }
//	    return option.Some(func(db Database) string {
//	        return fmt.Sprintf("Query on %s with timeout %d",
//	            db.ConnectionString, cfg.Timeout)
//	    })
//	}
//
//	// Sequenced: takes Database first, then Config
//	sequenced := SequenceReader(original)
//
//	db := Database{ConnectionString: "localhost:5432"}
//	cfg := Config{Timeout: 30}
//
//	// Apply database first to get a function that takes config
//	configReader := sequenced(db)
//	// Then apply config to get the final result
//	result := configReader(cfg)
//	// result is Option[string]
func SequenceReader[R1, R2, A any](ma ReaderOption[R2, Reader[R1, A]]) reader.Kleisli[R2, R1, Option[A]] {
	return readert.SequenceReader(
		option.Map,
		ma,
	)
}

func Traverse[R2, R1, A, B any](
	f Kleisli[R1, A, B],
) func(ReaderOption[R2, A]) Kleisli[R2, R1, B] {
	return readert.Traverse[ReaderOption[R2, A]](
		option.Map,
		option.Chain,
		f,
	)
}

func TraverseReader[R2, R1, A, B any](
	f reader.Kleisli[R1, A, B],
) func(ReaderOption[R2, A]) Kleisli[R2, R1, B] {
	return readert.TraverseReader[ReaderOption[R2, A]](
		option.Map,
		option.Map,
		f,
	)
}
