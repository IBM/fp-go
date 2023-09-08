// Copyright (c) 2023 IBM Corp.
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
	"fmt"

	A "github.com/IBM/fp-go/array"
	E "github.com/IBM/fp-go/either"
	"github.com/IBM/fp-go/errors"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	O "github.com/IBM/fp-go/option"
	P "github.com/IBM/fp-go/predicate"
	S "github.com/IBM/fp-go/string"
)

var (
	// httpGet :: Route -> Task Error JSON
	httpGet = F.Flow2(
		S.Format[string]("json for %s"),
		IOE.Of[error, string],
	)

	// routes :: Map Route Route
	routes = map[string]string{
		"/":      "/",
		"/about": "/about",
	}

	// validate :: Player -> Either error Player
	validatePlayer = E.FromPredicate(P.ContraMap(Player.getName)(S.IsNonEmpty), F.Flow2(Player.getId, errors.OnSome[int]("player %d must have a name")))

	// readfile :: String -> String -> Task Error String
	readfile = F.Curry2(func(encoding, file string) IOE.IOEither[error, string] {
		return IOE.Of[error](fmt.Sprintf("content of %s (%s)", file, encoding))
	})

	//   readdir :: String -> Task Error [String]
	readdir = IOE.Of[error](A.From("file1", "file2", "file3"))
)

func Example_solution12A() {
	// getJsons :: Map Route Route -> Task Error (Map Route JSON)
	getJsons := IOE.TraverseRecord[string](httpGet)

	fmt.Println(getJsons(routes)())

	// Output:
	// Right[<nil>, map[string]string](map[/:json for / /about:json for /about])
}

func Example_solution12B() {
	// startGame :: [Player] -> [Either Error String]
	startGame := F.Flow2(
		E.TraverseArray(validatePlayer),
		E.MapTo[error, []Player]("Game started"),
	)

	fmt.Println(startGame(A.From(playerAlbert, playerTheresa)))
	fmt.Println(startGame(A.From(playerAlbert, Player{Id: 4})))

	// Output:
	// Right[<nil>, string](Game started)
	// Left[*errors.errorString, string](player 4 must have a name)
}

func Example_solution12C() {
	traverseO := O.Traverse[string](
		IOE.Of[error, O.Option[string]],
		IOE.Map[error, string, O.Option[string]],
	)

	// readFirst :: String -> Task Error (Maybe String)
	readFirst := F.Pipe2(
		readdir,
		IOE.Map[error](A.Head[string]),
		IOE.Chain(traverseO(readfile("utf-8"))),
	)

	fmt.Println(readFirst())

	// Output:
	// Right[<nil>, option.Option[string]](Some[string](content of file1 (utf-8)))
}
