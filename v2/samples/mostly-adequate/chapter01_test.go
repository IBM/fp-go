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

package mostlyadequate

import "fmt"

type Flock struct {
	Seagulls int
}

func MakeFlock(n int) Flock {
	return Flock{Seagulls: n}
}

func (f *Flock) Conjoin(other *Flock) *Flock {
	f.Seagulls += other.Seagulls
	return f
}

func (f *Flock) Breed(other *Flock) *Flock {
	f.Seagulls *= other.Seagulls
	return f
}

func Example_flock() {

	flockA := MakeFlock(4)
	flockB := MakeFlock(2)
	flockC := MakeFlock(0)

	fmt.Println(flockA.Conjoin(&flockC).Breed(&flockB).Conjoin(flockA.Breed(&flockB)).Seagulls)

	// Output: 32
}
