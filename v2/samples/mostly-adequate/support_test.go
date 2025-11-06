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

import (
	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
)

type (
	Car struct {
		Name        string
		Horsepower  int
		DollarValue float32
		InStock     bool
	}
)

func (car Car) getInStock() bool {
	return car.InStock
}

func (car Car) getDollarValue() float32 {
	return car.DollarValue
}

func (car Car) getHorsepower() int {
	return car.Horsepower
}

func (car Car) getName() string {
	return car.Name
}

func average(val []float32) float32 {
	return F.Pipe2(
		val,
		A.Fold(N.MonoidSum[float32]()),
		N.Div(float32(len(val))),
	)
}

var (
	Cars = A.From(Car{
		Name:        "Ferrari FF",
		Horsepower:  660,
		DollarValue: 700000,
		InStock:     true,
	}, Car{
		Name:        "Spyker C12 Zagato",
		Horsepower:  650,
		DollarValue: 648000,
		InStock:     false,
	}, Car{
		Name:        "Jaguar XKR-S",
		Horsepower:  550,
		DollarValue: 132000,
		InStock:     true,
	}, Car{
		Name:        "Audi R8",
		Horsepower:  525,
		DollarValue: 114200,
		InStock:     false,
	}, Car{
		Name:        "Aston Martin One-77",
		Horsepower:  750,
		DollarValue: 1850000,
		InStock:     true,
	}, Car{
		Name:        "Pagani Huayra",
		Horsepower:  700,
		DollarValue: 1300000,
		InStock:     false,
	})
)
