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

// Package generic provides FromOptional and FromArrayOptional, which turn optionals into
// traversal endomorphisms. FromOptional focuses on zero or one value, FromArrayOptional on
// the elements of an array that may be absent. Where the optional does not match, the
// traversal yields the identity endomorphism.
package generic
