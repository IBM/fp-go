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

package examples

import (
	"fmt"
	"strconv"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	S "github.com/IBM/fp-go/v2/string"
)

func validatePort(port int) (int, error) {
	if port > 0 {
		return port, nil
	}
	return 0, fmt.Errorf("Value %d is not a valid port number", port)
}

func Example_either_monad() {

	// func(string) E.Either[error, int]
	atoi := E.Eitherize1(strconv.Atoi)
	// func(int) E.Either[error, int]
	valPort := E.Eitherize1(validatePort)

	// func(string) E.Either[error, string]
	makeUrl := F.Flow3(
		atoi,
		E.Chain(valPort),
		E.Map[error](S.Format[int]("http://localhost:%d")),
	)

	fmt.Println(makeUrl("8080"))

	// Output:
	// Right[string](http://localhost:8080)
}

func Example_either_idiomatic() {

	makeUrl := func(port string) (string, error) {
		parsed, err := strconv.Atoi(port)
		if err != nil {
			return "", err
		}
		valid, err := validatePort(parsed)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("http://localhost:%d", valid), nil
	}

	url, err := makeUrl("8080")
	if err != nil {
		panic(err)
	}
	fmt.Println(url)

	// Output:
	// http://localhost:8080
}

func Example_either_worlds() {

	// func(string) E.Either[error, int]
	atoi := E.Eitherize1(strconv.Atoi)
	// func(int) E.Either[error, int]
	valPort := E.Eitherize1(validatePort)

	// func(string) E.Either[error, string]
	makeUrl := F.Flow3(
		atoi,
		E.Chain(valPort),
		E.Map[error](S.Format[int]("http://localhost:%d")),
	)

	// func(string) (string, error)
	makeUrlGo := E.Uneitherize1(makeUrl)

	url, err := makeUrlGo("8080")
	if err != nil {
		panic(err)
	}
	fmt.Println(url)

	// Output:
	// http://localhost:8080
}
