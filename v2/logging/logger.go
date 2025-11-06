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

package logging

import (
	"log"
)

func LoggingCallbacks(loggers ...*log.Logger) (func(string, ...any), func(string, ...any)) {
	switch len(loggers) {
	case 0:
		def := log.Default()
		return def.Printf, def.Printf
	case 1:
		log0 := loggers[0]
		return log0.Printf, log0.Printf
	default:
		return loggers[0].Printf, loggers[1].Printf
	}
}
